package model

type Interaction struct {
	LikeCount    int `db:"like_count"`
	RetweetCount int `db:"retweet_count"`
	ReplyCount   int `db:"reply_count"`
}
