package model

type LikeCount struct {
	TweetID int `db:"tweet_id"`
	Count   int `db:"count"`
}
