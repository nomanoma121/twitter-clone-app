package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/model"
	"time"

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
	g.GET("/tweets/all", h.GetAllTweets)
	g.POST("/tweet", h.CreateTweet)
	g.POST("/retweet", h.Retweet)
}

type GetAllTweetsResponseUser struct {
	ID    int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name  string `json:"name"`
}

type GetAllTweetsResponseRetweet struct {
	ID      int                      `json:"id"`
	User    GetAllTweetsResponseUser `json:"user"`
	Content string                   `json:"content"`
	CreatedAt time.Time               `json:"created_at"`
}

type GetAllTweetsResponseReply struct {
	ID      int                      `json:"id"`
	User    GetAllTweetsResponseUser `json:"user"`
	Content string                   `json:"content"`
	CreatedAt time.Time               `json:"created_at"`
}

type GetAllTweetsResponse struct {
	ID        int                          `json:"id"`
	User      GetAllTweetsResponseUser     `json:"user"`
	Content   string                       `json:"content"`
	Retweet   *GetAllTweetsResponseRetweet `json:"retweet"`
	Reply     *GetAllTweetsResponseReply   `json:"reply"`
	CreatedAt time.Time                    `json:"created_at"`
}

// Retweet or Replyがある場合はそれを取得し、Tweetに加える関数
func getRetweetOrReply(actions *[]model.Tweet, tweets *[]model.Tweet) {
	// 参照先のTweetのIDを格納するスライス
	var actionsIDs []int
	for _, tweet := range tweets {
		// RetweetとReplyを同時に行うことはない
		if tweet.RetweetID != nil {
			actionsIDs = append(actionsIDs, *tweet.RetweetID)
		} else if tweet.ReplyID != nil {
			actionsIDs = append(actionsIDs, *tweet.ReplyID)
		}
	}
	if len(actionsIDs) > 0 {
		// retweetかreplyする参照先のTweetをしたUser情報を取得
		query, args, err := sqlx.In(`
			SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.display_id as "user.display_id"
			FROM tweets
			JOIN users ON tweets.user_id = users.id
			WHERE tweets.id IN (?)
		`, actionsIDs)
		if err != nil {
			log.Println(err)
		}

	err = h.db.Select(actions, query, args...)
	if err != nil {
		log.Println(err)
	}
	// retweetかreplyが参照するTweetをretweetかreplyに加える
	var actionsMap = map[int]model.Tweet{}
	for _, action := range *actions {
		actionsMap[action.ID] = action
	}
	// retweetかreplyが参照するTweetをretweetかreplyに加える
	for i, tweet := range *tweets {
		if tweet.RetweetID != nil {
			action, ok := actionsMap[*tweet.RetweetID]
			if !ok {
				return
			}
			(*tweets)[i].Retweet = &action
		} else if tweet.ReplyID != nil {
			action, ok := actionsMap[*tweet.ReplyID]
			if !ok {
				return
			}
			(*tweets)[i].Reply = &action
		}
	}
}

func (h *TweetHandler) GetAllTweets(c echo.Context) error {
	var tweets []model.Tweet
	// データベースからすべてのツイートを取得
	err := h.db.Select(&tweets, `
		SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.display_id as "user.display_id"
		FROM tweets
		JOIN users ON tweets.user_id = users.id
	`)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetAllTweetsResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	var retweets []model.Tweet
	var replies []model.Tweet
	// RetweetとReplyを取得
	getRetweetOrReply(&retweets, &tweets)
	getRetweetOrReply(&replies, &tweets)

	res := make([]GetAllTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetAllTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetAllTweetsResponseRetweet{
				ID:      tweet.Retweet.ID,
				User:    GetAllTweetsResponseUser{ID: tweet.Retweet.User.ID, Name: tweet.Retweet.User.Name, DisplayID: tweet.Retweet.User.DisplayID},
				Content: tweet.Retweet.Content,
				CreatedAt: tweet.Retweet.CreatedAt,
			}
		}
		reply := (*GetAllTweetsResponseReply)(nil)
		if tweet.Reply != nil {
			reply = &GetAllTweetsResponseReply{
				ID:      tweet.Reply.ID,
				User: 	GetAllTweetsResponseUser{ID: tweet.Reply.User.ID, Name: tweet.Reply.User.Name, DisplayID: tweet.Reply.User.DisplayID},
				Content: tweet.Reply.Content,
				CreatedAt: tweet.Reply.CreatedAt,
			}
		}
		res[i] = GetAllTweetsResponse{
			ID: tweet.ID,
			User: GetAllTweetsResponseUser{
				ID:    tweet.User.ID,
				Name:  tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
			},
			Content: tweet.Content,
			Retweet: retweet,
			Reply: reply,
			CreatedAt: tweet.CreatedAt,
		}
	}

	return c.JSON(200, res)
}


type GetTweetsResponseUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetTweetsResponseRetweet struct {
	ID      int                   `json:"id"`
	User    GetTweetsResponseUser `json:"user"`
	Content string                `json:"content"`
}

type GetTweetsResponse struct {
	ID      int                       `json:"id"`
	User    GetTweetsResponseUser     `json:"user"`
	Content string                    `json:"content"`
	Retweet *GetTweetsResponseRetweet `json:"retweet"`
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
	var retweets []model.Tweet
	var retweetIDs []int
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetIDs = append(retweetIDs, *tweet.RetweetID)
		}
	}
	if len(retweetIDs) > 0 {
		query, args, err := sqlx.In(`
			SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.email as "user.email"
			FROM tweets
			JOIN users ON tweets.user_id = users.id
			WHERE tweets.id IN (?)
		`, retweetIDs)
		fmt.Printf("query: %v\n", query)
		fmt.Printf("args: %v\n", args)
		if err != nil {
			log.Println(err)
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
		err = h.db.Select(&retweets, query, args...)
		if err != nil {
			log.Println(err)
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
	}
	var retweetMap = map[int]model.Tweet{}
	for _, retweet := range retweets {
		retweetMap[retweet.ID] = retweet
	}
	for i, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweet, ok := retweetMap[*tweet.RetweetID]
			if !ok {
				return c.JSON(500, map[string]string{"message": "Internal Server Error"})
			}
			tweets[i].Retweet = &retweet
		}
	}

	fmt.Printf("retweets: %#v\n", retweets)

	for _, tweet := range tweets {
		fmt.Printf("tweet: %#v\n", tweet)
	}
	fmt.Printf("tweets len: %#v\n", len(tweets))

	res := make([]GetTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetTweetsResponseRetweet{
				ID:      tweet.Retweet.ID,
				User:    GetTweetsResponseUser{ID: tweet.Retweet.User.ID, Name: tweet.Retweet.User.Name, Email: tweet.Retweet.User.Email},
				Content: tweet.Retweet.Content,
			}
		}
		res[i] = GetTweetsResponse{
			ID: tweet.ID,
			User: GetTweetsResponseUser{
				ID:    tweet.User.ID,
				Name:  tweet.User.Name,
				Email: tweet.User.Email,
			},
			Content: tweet.Content,
			Retweet: retweet,
		}
	}

	return c.JSON(200, res)
}

type CreateTweetRequest struct {
	Content string `json:"content" validate:"required"`
}

func (h *TweetHandler) CreateTweet(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(CreateTweetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	fmt.Printf("req: %#v\n", req)

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	_, err := h.db.Exec("INSERT INTO tweets (user_id, content) VALUES (?, ?)", userID, req.Content)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(201)
}

type RetweetRequest struct {
	TweetID int     `json:"tweet_id" validate:"required"`
	Content *string `json:"content"`
}

func (h *TweetHandler) Retweet(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(RetweetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	var tweet model.Tweet
	err := h.db.Get(&tweet, "SELECT * FROM tweets WHERE id = ?", req.TweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "Not Found"})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	if tweet.RetweetID != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if req.Content == nil {
		_, err = h.db.Exec("INSERT INTO tweets (user_id, retweet_id) VALUES (?, ?)", userID, req.TweetID)
		if err != nil {
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
	} else {
		_, err = h.db.Exec("INSERT INTO tweets (user_id, retweet_id, content) VALUES (?, ?, ?)", userID, req.TweetID, *req.Content)
		if err != nil {
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
	}

	return c.NoContent(201)
}
