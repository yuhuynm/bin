package entity

import "time"

type Secret struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Token            string     `gorm:"size:64;uniqueIndex;not null" json:"token"`
	EncryptedContent string     `gorm:"type:text;not null" json:"-"`
	ContentType      string     `gorm:"size:20;not null;default:text" json:"content_type"`
	HasPassword      bool       `gorm:"not null;default:false" json:"has_password"`
	PasswordHash     *string    `gorm:"type:text" json:"-"`
	ExpiresAt        *time.Time `json:"expires_at"`
	MaxViews         int        `gorm:"not null;default:0" json:"max_views"`
	ViewCount        int        `gorm:"not null;default:0" json:"view_count"`
	IsConsumed       bool       `gorm:"not null;default:false" json:"is_consumed"`
	ConsumedAt       *time.Time `json:"consumed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (Secret) TableName() string {
	return "Secrets"
}
