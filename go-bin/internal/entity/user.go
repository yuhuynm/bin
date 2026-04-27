package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"not null" json:"tenant_id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Email     string    `gorm:"size:160;uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "Users"
}
