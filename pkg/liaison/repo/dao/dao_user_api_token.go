package dao

import (
	"time"

	"github.com/liaisonio/liaison/pkg/liaison/repo/model"
)

func (d *dao) CreateUserAPIToken(tok *model.UserAPIToken) error {
	return d.getDB().Create(tok).Error
}

// ListUserAPITokens returns all non-revoked tokens for a user, newest first.
func (d *dao) ListUserAPITokens(userID uint) ([]*model.UserAPIToken, error) {
	var out []*model.UserAPIToken
	err := d.getDB().
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Order("id DESC").
		Find(&out).Error
	return out, err
}

// GetUserAPITokensByPrefix returns non-revoked tokens whose prefix matches.
// The caller bcrypt-compares each candidate against the full plaintext — a
// prefix collision is extremely rare but theoretically possible.
func (d *dao) GetUserAPITokensByPrefix(prefix string) ([]*model.UserAPIToken, error) {
	var out []*model.UserAPIToken
	err := d.getDB().
		Where("token_prefix = ? AND revoked_at IS NULL", prefix).
		Find(&out).Error
	return out, err
}

// RevokeUserAPIToken only succeeds if the token belongs to userID — prevents
// cross-user revocation via id guessing.
func (d *dao) RevokeUserAPIToken(userID, id uint) error {
	now := time.Now()
	return d.getDB().
		Model(&model.UserAPIToken{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", id, userID).
		Update("revoked_at", &now).Error
}

func (d *dao) TouchUserAPIToken(id uint, ip string) error {
	now := time.Now()
	return d.getDB().
		Model(&model.UserAPIToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_used_at": &now,
			"last_used_ip": ip,
		}).Error
}

func (d *dao) CountUserAPITokens(userID uint) (int64, error) {
	var n int64
	err := d.getDB().
		Model(&model.UserAPIToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Count(&n).Error
	return n, err
}
