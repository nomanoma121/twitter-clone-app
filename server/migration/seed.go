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
		Email:    "john.doe@example.com",
		Password: "password123",
	},
	{
		Email:    "jane.smith@example.com",
		Password: "securepass",
	},
	{
		Email:    "alice.brown@example.com",
		Password: "alice2023",
	},
	{
		Email:    "charlie.wilson@example.com",
		Password: "wilson42",
	},
	{
		Email:    "emily.jones@example.com",
		Password: "emilyPass1",
	},
	{
		Email:    "michael.johnson@example.com",
		Password: "mikeSecure",
	},
	{
		Email:    "olivia.taylor@example.com",
		Password: "oliviaTay99",
	},
	{
		Email:    "william.miller@example.com",
		Password: "millerRocks",
	},
	{
		Email:    "sophia.davis@example.com",
		Password: "sophiaD",
	},
	{
		Email:    "james.moore@example.com",
		Password: "mooreJames",
	},
	{
		Email:    "isabella.jackson@example.com",
		Password: "jacksonBella",
	},
	{
		Email:    "ethan.thomas@example.com",
		Password: "thomasEthan",
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
		Name:      "John Doe",
		DisplayID: "john_d",
		IconURL:   "http://localhost:5173/public/images/icon/icon_1.webp",
		HeaderURL: "http://localhost:5173/public/images/header/header_1.jpg",
		Profile:   "Loves coding and exploring new technologies.",
	},
	{
		UserID:    2,
		Name:      "Jane Smith",
		DisplayID: "jane_s",
		IconURL:   "http://localhost:5173/public/images/icon/icon_2.webp",
		HeaderURL: "http://localhost:5173/public/images/header/header_2.jpg",
		Profile:   "Graphic designer with a passion for art.",
	},
	{
		UserID:    3,
		Name:      "Alice Brown",
		DisplayID: "alice_b",
		IconURL:   "http://localhost:5173/public/images/icon/icon_3.webp",
		HeaderURL: "http://localhost:5173/public/images/header/header_3.jpg",
		Profile:   "Bookworm and coffee enthusiast.",
	},
	{
		UserID:    4,
		Name:      "Charlie Wilson",
		DisplayID: "charlie_w",
		IconURL:   "http://localhost:5173/public/images/icon/icon_4.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_4.jpg",
		Profile:   "Travel blogger sharing adventures.",
	},
	{
		UserID:    5,
		Name:      "Emily Jones",
		DisplayID: "emily_j",
		IconURL:   "http://localhost:5173/public/images/icon/icon_5.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_5.jpg",
		Profile:   "Photographer capturing moments.",
	},
	{
		UserID:    6,
		Name:      "Michael Johnson",
		DisplayID: "mike_j",
		IconURL:   "http://localhost:5173/public/images/icon/icon_6.webp",
		HeaderURL: "http://localhost:5173/public/images/header/header_6.jpg",
		Profile:   "Tech geek and gamer.",
	},
	{
		UserID:    7,
		Name:      "Olivia Taylor",
		DisplayID: "olivia_t",
		IconURL:   "http://localhost:5173/public/images/icon/icon_7.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_7.jpg",
		Profile:   "Fashion designer and trendsetter.",
	},
	{
		UserID:    8,
		Name:      "William Miller",
		DisplayID: "will_m",
		IconURL:   "http://localhost:5173/public/images/icon/icon_8.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_8.jpg",
		Profile:   "Fitness coach helping others achieve goals.",
	},
	{
		UserID:    9,
		Name:      "Sophia Davis",
		DisplayID: "sophia_d",
		IconURL:   "http://localhost:5173/public/images/icon/icon_9.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_9.jpg",
		Profile:   "Musician composing melodies.",
	},
	{
		UserID:    10,
		Name:      "James Moore",
		DisplayID: "james_m",
		IconURL:   "http://localhost:5173/public/images/icon/icon_10.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_10.jpg",
		Profile:   "Entrepreneur building startups.",
	},
	{
		UserID:    11,
		Name:      "Isabella Jackson",
		DisplayID: "bella_j",
		IconURL:   "http://localhost:5173/public/images/icon/icon_11.webp",
		HeaderURL: "http://localhost:5173/public/images/header/header_11.jpg",
		Profile:   "Chef creating unique recipes.",
	},
	{
		UserID:    12,
		Name:      "Ethan Thomas",
		DisplayID: "ethan_t",
		IconURL:   "http://localhost:5173/public/images/icon/icon_12.jpg",
		HeaderURL: "http://localhost:5173/public/images/header/header_12.jpg",
		Profile:   "Software developer specializing in AI.",
	},
}

type InsertTweet struct {
	UserID  int
	Content string
}

// var tweets = []InsertTweet{
// 	{
// 		UserID:  1,
// 		Content: "tweet1",
// 	},
// 	{
// 		UserID:  2,
// 		Content: "tweet2",
// 	},
// }
// 日本語でツイートのテストデータを20個作成してください。実際のツイートのような内容にしてほしいです。
// 例：「今日はいい天気ですね。」

var tweets = []InsertTweet{
	{
		UserID:  1,
		Content: "今日はいい天気ですね。",
	},
	{
		UserID:  2,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  3,
		Content: "今日は暑いです。",
	},
	{
		UserID:  4,
		Content: "今日は寒いです。",
	},
	{
		UserID:  5,
		Content: "今日は曇りです。",
	},
	{
		UserID:  6,
		Content: "今日は雪が降っています。",
	},
	{
		UserID:  7,
		Content: "今日は風が強いです。",
	},
	{
		UserID:  8,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  9,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  10,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  11,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  12,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  1,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  2,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  3,
		Content: "今日は雨が降っています。",
	},
	{
		UserID:  4,
		Content: "今日は雨が降っています。",
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
		URL:     "http://localhost:5173/public/default-icon.jpeg",
	},
	{
		TweetID: 2,
		URL:     "http://localhost:5173/public/default-icon.jpeg",
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
