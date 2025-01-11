package model

type CountResult struct {
	TweetID int `db:"tweet_id"`
	Count   int `db:"count"`
}
