package model

import (
	"time"
)

type Tweet struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	User      User      `db:"user"`
	Content   string    `db:"content"`
	RetweetID *int      `db:"retweet_id"`
	Retweet   *Tweet    `db:"retweet"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
