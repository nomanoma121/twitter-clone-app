package model

import (
	"time"
)

type TweetImage struct {
	ID        int       `db:"id"`
	TweetID   int       `db:"tweet_id"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
