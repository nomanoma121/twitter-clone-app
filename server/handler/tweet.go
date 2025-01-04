package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/model"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type TweetHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
}

func NewTweetHandler(db *sqlx.DB) *TweetHandler {
	return &TweetHandler{db: db, validator: validator.New()}
}

func (h *TweetHandler) Register(g *echo.Group) {
	g.GET("/tweets", h.GetTweets)
}

type GetTweetsResponseUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetTweetsResponse struct {
	ID        int                   `json:"id"`
	User      GetTweetsResponseUser `json:"user"`
	Content   string                `json:"content"`
	RetweetID *int                  `json:"retweet_id"`
}

func (h *TweetHandler) GetTweets(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var tweets []model.Tweet
	// err := h.db.Select(&tweets, "SELECT * FROM tweet WHERE user_id = ?", userID)
	// 自分のツイートと一緒にユーザーをJoinして取得
	err := h.db.Select(&tweets, `
		SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.email as "user.email"
		FROM tweets
		JOIN users ON tweets.user_id = users.id
		WHERE tweets.user_id = ?
	`, userID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetTweetsResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	// fmt.Printf("tweets: %#v\n", tweets)
	for _, tweet := range tweets {
		fmt.Printf("tweet: %#v\n", tweet)
	}
	fmt.Printf("tweets len: %#v\n", len(tweets))

	res := make([]GetTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		res[i] = GetTweetsResponse{
			ID: tweet.ID,
			User: GetTweetsResponseUser{
				ID:    tweet.User.ID,
				Name:  tweet.User.Name,
				Email: tweet.User.Email,
			},
			Content:   tweet.Content,
			RetweetID: tweet.RetweetID,
		}
	}

	return c.JSON(200, res)
}
