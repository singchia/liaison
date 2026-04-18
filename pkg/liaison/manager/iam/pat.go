package iam

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/liaisonio/liaison/pkg/liaison/repo/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	patRandomBytes       = 25 // 25 bytes = 200 bits = 40 base32 chars
	patPrefixRandomChars = 8  // chars kept after "liaison_pat_" for the stored prefix
	patPrefixTotalLen    = len(model.PATPlaintextPrefix) + patPrefixRandomChars
	patMaxPerUser        = 50
	patMaxNameLen        = 64
	patBcryptCost        = 10
)

// Sentinel errors so HTTP handlers can map to status codes without string-matching.
var (
	ErrPATNameRequired = errors.New("token name is required")
	ErrPATNameTooLong  = fmt.Errorf("token name must be at most %d characters", patMaxNameLen)
	ErrPATTooMany      = fmt.Errorf("token limit reached (max %d per user)", patMaxPerUser)
	ErrPATNotFound     = errors.New("token not found")
	ErrPATInvalid      = errors.New("invalid or revoked token")
	ErrPATExpired      = errors.New("token expired")
)

// PATCreateResult carries the plaintext token back to the caller exactly once.
// The Token field is the only time plaintext is ever exposed.
type PATCreateResult struct {
	Model *model.UserAPIToken
	Token string // plaintext, shown to user exactly once
}

// CreatePAT mints a new personal access token for userID.
// expiresAt == nil means the token never expires.
func (s *IAMService) CreatePAT(userID uint, name string, expiresAt *time.Time) (*PATCreateResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrPATNameRequired
	}
	if len(name) > patMaxNameLen {
		return nil, ErrPATNameTooLong
	}

	count, err := s.repo.CountUserAPITokens(userID)
	if err != nil {
		return nil, err
	}
	if count >= patMaxPerUser {
		return nil, ErrPATTooMany
	}

	raw := make([]byte, patRandomBytes)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("generate random token: %w", err)
	}
	suffix := strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString(raw), "="))
	plaintext := model.PATPlaintextPrefix + suffix

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), patBcryptCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt hash: %w", err)
	}

	prefix := plaintext
	if len(prefix) > patPrefixTotalLen {
		prefix = prefix[:patPrefixTotalLen]
	}

	tok := &model.UserAPIToken{
		UserID:      userID,
		Name:        name,
		TokenPrefix: prefix,
		TokenHash:   string(hash),
		Scope:       "full_access",
		ExpiresAt:   expiresAt,
	}
	if err := s.repo.CreateUserAPIToken(tok); err != nil {
		return nil, err
	}
	return &PATCreateResult{Model: tok, Token: plaintext}, nil
}

// ListPATs returns the user's non-revoked tokens, with hashes stripped so a
// caller cannot accidentally leak them via the HTTP response.
func (s *IAMService) ListPATs(userID uint) ([]*model.UserAPIToken, error) {
	out, err := s.repo.ListUserAPITokens(userID)
	if err != nil {
		return nil, err
	}
	for _, t := range out {
		t.TokenHash = ""
	}
	return out, nil
}

// RevokePAT is idempotent; an unknown id for the caller returns nil (same as
// already-revoked) to avoid leaking ownership info.
func (s *IAMService) RevokePAT(userID, id uint) error {
	return s.repo.RevokeUserAPIToken(userID, id)
}

// GetUserByPAT authenticates a bearer token as a PAT and returns the owning user.
// Called from the auth middleware when the token starts with "liaison_pat_".
func (s *IAMService) GetUserByPAT(plaintext, requestIP string) (*model.User, error) {
	plaintext = strings.TrimSpace(plaintext)
	if !strings.HasPrefix(plaintext, model.PATPlaintextPrefix) {
		return nil, ErrPATInvalid
	}
	prefix := plaintext
	if len(prefix) > patPrefixTotalLen {
		prefix = prefix[:patPrefixTotalLen]
	}

	candidates, err := s.repo.GetUserAPITokensByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	var match *model.UserAPIToken
	for _, c := range candidates {
		if err := bcrypt.CompareHashAndPassword([]byte(c.TokenHash), []byte(plaintext)); err == nil {
			match = c
			break
		}
	}
	if match == nil {
		return nil, ErrPATInvalid
	}
	if match.ExpiresAt != nil && time.Now().After(*match.ExpiresAt) {
		return nil, ErrPATExpired
	}

	user, err := s.repo.GetUserByID(match.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrPATInvalid
	}

	// Best-effort last_used update; don't fail auth if DB write fails.
	go func(id uint, ip string) {
		if err := s.repo.TouchUserAPIToken(id, ip); err != nil {
			log.Warnf("touch user_api_token id=%d: %v", id, err)
		}
	}(match.ID, requestIP)

	return user, nil
}
