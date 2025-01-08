package migration

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type InsertUser struct {
	Email    string
	Password string
}

var users = []InsertUser{
	{
		Email:    "user1@example.com",
		Password: "password1",
	},
	{
		Email:    "user2@example.com",
		Password: "password2",
	},
	{
		Email:    "user3@example.com",
		Password: "password3",
	},
	{
		Email:    "user4@example.com",
		Password: "password4",
	},
}

type InsertUserProfile struct {
	UserID    int
	Name      string
	DisplayID string
	IconURL   string
	HeaderURL string
	Profile   string
}

var userProfiles = []InsertUserProfile{
	{
		UserID:    1,
		Name:      "user1",
		DisplayID: "user1",
		IconURL:   "https://example.com/icon1.png",
		HeaderURL: "https://example.com/header1.png",
		Profile:   "profile1",
	},
	{
		UserID:    2,
		Name:      "user2",
		DisplayID: "user2",
		IconURL:   "https://example.com/icon2.png",
		HeaderURL: "https://example.com/header2.png",
		Profile:   "profile2",
	},
	{
		UserID:    3,
		Name:      "user3",
		DisplayID: "user3",
		IconURL:   "https://example.com/icon3.png",
		HeaderURL: "https://example.com/header3.png",
		Profile:   "profile3",
	},
	{
		UserID:    4,
		Name:      "user4",
		DisplayID: "user4",
		IconURL:   "https://example.com/icon4.png",
		HeaderURL: "https://example.com/header4.png",
		Profile:   "profile4",
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

type InsertReply struct {
	UserID  int
	Content string
	ReplyID int
}

var replies = []InsertReply{
	{
		UserID:  1,
		Content: "reply1",
		ReplyID: 2,
	},
	{
		UserID:  2,
		Content: "reply2",
		ReplyID: 1,
	},
}

type InsertTweetImage struct {
	TweetID int
	URL     string
}

var tweetImages = []InsertTweetImage{
	{
		TweetID: 1,
		URL:     "https://example.com/image1.png",
	},
	{
		TweetID: 2,
		URL:     "https://example.com/image2.png",
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

type InsertLike struct {
	user_id  int
	tweet_id int
}

var likes = []InsertLike{
	{
		user_id:  1,
		tweet_id: 2,
	},
	{
		user_id:  2,
		tweet_id: 1,
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
			INSERT INTO users (email, password_hash) VALUES (?, ?)
		`, user.Email, hash)
		if err != nil {
			return err
		}
	}

	for _, userProfile := range userProfiles {
		_, err := db.Exec(`
			INSERT INTO user_profiles (user_id, name, display_id, icon_url, header_url, profile) VALUES (?, ?, ?, ?, ?, ?)
		`, userProfile.UserID, userProfile.Name, userProfile.DisplayID, userProfile.IconURL, userProfile.HeaderURL, userProfile.Profile)
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

	for _, reply := range replies {
		_, err := db.Exec(`
			INSERT INTO tweets (user_id, content, reply_id) VALUES (?, ?, ?)
		`, reply.UserID, reply.Content, reply.ReplyID)
		if err != nil {
			return err
		}
	}

	for _, tweetImage := range tweetImages {
		_, err := db.Exec(`
			INSERT INTO tweet_images (tweet_id, url) VALUES (?, ?)
		`, tweetImage.TweetID, tweetImage.URL)
		if err != nil {
			return err
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

	for _, like := range likes {
		_, err := db.Exec(`
			INSERT INTO likes (user_id, tweet_id) VALUES (?, ?)
		`, like.user_id, like.tweet_id)
		if err != nil {
			return err
		}
	}

	return nil
}
