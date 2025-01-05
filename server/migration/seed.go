package migration

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type InsertUser struct {
	Name     string
	Email    string
	Password string
}

var users = []InsertUser{
	{
		Name:     "user1",
		Email:    "user1@example.com",
		Password: "password1",
	},
	{
		Name:     "user2",
		Email:    "user2@example.com",
		Password: "password2",
	},
	{
		Name:     "user3",
		Email:    "user3@example.com",
		Password: "password3",
	},
	{
		Name:     "user4",
		Email:    "user4@example.com",
		Password: "password4",
	},
}

type InsertTweet struct {
	UserID  int
	Content string
}

var tweets = []InsertTweet{
	{
		UserID:  1,
		Content: "tweet1",
	},
	{
		UserID:  2,
		Content: "tweet2",
	},
}

type InsertRetweet struct {
	UserID    int
	Content   *string
	RetweetID int
}

var retweetContent = "retweet1"

var retweets = []InsertRetweet{
	{
		UserID:    1,
		Content:   nil,
		RetweetID: 2,
	},
	{
		UserID:    2,
		Content:   &retweetContent,
		RetweetID: 1,
	},
}

type InsertFollow struct {
	FollowerID int
	FolloweeID int
}

var follows = []InsertFollow{
	{
		FollowerID: 1,
		FolloweeID: 2,
	},
	{
		FollowerID: 2,
		FolloweeID: 1,
	},
	{
		FollowerID: 2,
		FolloweeID: 4,
	},
	{
		FollowerID: 3,
		FolloweeID: 2,
	},
}

func Seed(db *sqlx.DB) error {
	err := Reset(db)
	if err != nil {
		return err
	}

	for _, user := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)
		`, user.Name, user.Email, hash)
		if err != nil {
			return err
		}
	}

	for _, tweet := range tweets {
		_, err := db.Exec(`
			INSERT INTO tweets (user_id, content) VALUES (?, ?)
		`, tweet.UserID, tweet.Content)
		if err != nil {
			return err
		}
	}

	for _, retweet := range retweets {
		if retweet.Content == nil {
			_, err := db.Exec(`
				INSERT INTO tweets (user_id, content, retweet_id) VALUES (?, ?, ?)
			`, retweet.UserID, "", retweet.RetweetID)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(`
				INSERT INTO tweets (user_id, content, retweet_id) VALUES (?, ?, ?)
			`, retweet.UserID, *retweet.Content, retweet.RetweetID)
			if err != nil {
				return err
			}
		}
	}

	for _, follow := range follows {
		_, err := db.Exec(`
			INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)
		`, follow.FollowerID, follow.FolloweeID)
		if err != nil {
			return err
		}
	}

	return nil
}
