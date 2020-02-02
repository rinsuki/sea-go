package models

import "database/sql"

type AccessToken struct {
	ID        uint         `gorm:"id"`
	UserID    uint         `gorm:"user_id"`
	AppID     uint         `gorm:"application_id"`
	Token     string       `gorm:"token"`
	RevokedAt sql.NullTime `gorm:"revoked_at"`
}
