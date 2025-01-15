package model

import (
	"time"
)

type UserProfile struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Name      string    `db:"name"`
	DisplayID string    `db:"display_id"`
	IconURL   string    `db:"icon_url"`
	HeaderURL string    `db:"header_url"`
	Profile   string    `db:"profile"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
