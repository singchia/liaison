package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/liaisonio/liaison/pkg/liaison/manager/iam"
	"github.com/liaisonio/liaison/pkg/liaison/manager/timefmt"
	"github.com/liaisonio/liaison/pkg/liaison/repo/model"
)

// ─── request/response types ──────────────────────────────────────────────────

type patCreateRequest struct {
	Name          string `json:"name"`
	ExpiresInDays int    `json:"expires_in_days,omitempty"` // 0 = never expires
}

type patView struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	TokenPrefix string  `json:"token_prefix"`
	Scope       string  `json:"scope"`
	LastUsedAt  *string `json:"last_used_at,omitempty"`
	LastUsedIP  string  `json:"last_used_ip,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type patCreateResponse struct {
	patView
	Token string `json:"token"` // plaintext, shown exactly once
}

func patViewFromModel(m *model.UserAPIToken) patView {
	v := patView{
		ID:          m.ID,
		Name:        m.Name,
		TokenPrefix: m.TokenPrefix,
		Scope:       m.Scope,
		LastUsedIP:  m.LastUsedIP,
		CreatedAt:   timefmt.FormatDateTime(m.CreatedAt),
	}
	if m.LastUsedAt != nil && !m.LastUsedAt.IsZero() {
		s := timefmt.FormatDateTime(*m.LastUsedAt)
		v.LastUsedAt = &s
	}
	if m.ExpiresAt != nil && !m.ExpiresAt.IsZero() {
		s := timefmt.FormatDateTime(*m.ExpiresAt)
		v.ExpiresAt = &s
	}
	return v
}

// ─── dispatchers ─────────────────────────────────────────────────────────────

// handleTokensHTTP multiplexes GET (list) and POST (create) on /api/v1/iam/tokens.
func (web *web) handleTokensHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		web.handleListTokens(w, r)
	case http.MethodPost:
		web.handleCreateToken(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"code":    http.StatusMethodNotAllowed,
			"message": "method not allowed",
		})
	}
}

// handleTokenByIDHTTP handles DELETE /api/v1/iam/tokens/{id}.
func (web *web) handleTokenByIDHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"code":    http.StatusMethodNotAllowed,
			"message": "method not allowed",
		})
		return
	}
	web.handleRevokeToken(w, r)
}

// ─── create ──────────────────────────────────────────────────────────────────

func (web *web) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}

	var req patCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"code":    http.StatusBadRequest,
			"message": "invalid request body",
		})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresInDays) * 24 * time.Hour)
		expiresAt = &t
	}

	result, err := web.iamService.CreatePAT(user.ID, req.Name, expiresAt)
	if err != nil {
		status, msg := patErrorStatus(err)
		writeJSON(w, status, map[string]any{"code": status, "message": msg})
		return
	}

	resp := patCreateResponse{
		patView: patViewFromModel(result.Model),
		Token:   result.Token,
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// ─── list ────────────────────────────────────────────────────────────────────

func (web *web) handleListTokens(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}

	tokens, err := web.iamService.ListPATs(user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	views := make([]patView, 0, len(tokens))
	for _, t := range tokens {
		views = append(views, patViewFromModel(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
		"data":    map[string]any{"tokens": views},
	})
}

// ─── revoke ──────────────────────────────────────────────────────────────────

func (web *web) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/iam/tokens/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"code":    http.StatusBadRequest,
			"message": "invalid token id",
		})
		return
	}

	if err := web.iamService.RevokePAT(user.ID, uint(id)); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
	})
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// patErrorStatus maps PAT service errors to HTTP status codes + user messages.
func patErrorStatus(err error) (int, string) {
	switch {
	case errors.Is(err, iam.ErrPATNameRequired), errors.Is(err, iam.ErrPATNameTooLong):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, iam.ErrPATTooMany):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
