package model

import (
	"time"
)

type Good struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TweetID   int       `db:"tweet_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
