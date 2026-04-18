package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// StringSlice stores []string as a JSON-encoded TEXT column.
// Empty slice serializes to "[]" (not NULL) so callers can distinguish
// "no entry" (allow-all) from "empty allowlist" (deny-all).
type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		*s = StringSlice{}
		return nil
	}
}

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// ProxyFirewallRule holds the source-IP allowlist for a single proxy.
// An absent row means "allow all"; an empty AllowedCIDRs slice means "deny all".
type ProxyFirewallRule struct {
	ID           uint `gorm:"primarykey;autoIncrement"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	UserID       uint           `gorm:"column:user_id;type:int;not null;default:0;index"`
	ProxyID      uint           `gorm:"column:proxy_id;type:int;not null;uniqueIndex"`
	AllowedCIDRs StringSlice    `gorm:"column:allowed_cidrs;type:text;not null"`
}

func (ProxyFirewallRule) TableName() string {
	return "proxy_firewall_rules"
}
