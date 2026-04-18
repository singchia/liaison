package model

import (
	"time"

	"gorm.io/gorm"
)

// UserAPIToken is a long-lived personal access token issued to a user for
// API/CLI use. Tokens are stored as bcrypt hashes; the plaintext is returned
// to the user exactly once at creation time.
//
// token_prefix stores a short, human-identifiable fragment of the plaintext so
// the UI can render "liaison_pat_a1b2..." and the auth middleware can narrow
// the hash candidates via an index lookup before bcrypt-comparing.
type UserAPIToken struct {
	gorm.Model
	UserID      uint       `gorm:"column:user_id;type:int;not null;index"`
	Name        string     `gorm:"column:name;type:varchar(64);not null"`
	TokenPrefix string     `gorm:"column:token_prefix;type:varchar(32);not null;index"`
	TokenHash   string     `gorm:"column:token_hash;type:varchar(128);not null"`
	Scope       string     `gorm:"column:scope;type:varchar(64);not null;default:'full_access'"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at"`
	LastUsedIP  string     `gorm:"column:last_used_ip;type:varchar(64);default:''"`
	ExpiresAt   *time.Time `gorm:"column:expires_at"` // nil = never expires
	RevokedAt   *time.Time `gorm:"column:revoked_at"` // non-nil = revoked
}

func (UserAPIToken) TableName() string {
	return "user_api_tokens"
}

// PATPlaintextPrefix is the fixed literal every PAT plaintext starts with.
// Used by the auth middleware to distinguish PATs from JWT session tokens.
const PATPlaintextPrefix = "liaison_pat_"
