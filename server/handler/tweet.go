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
	g.GET("/tweets/all", h.GetAllTweets)
	g.GET("/tweets/timeline", h.GetTweets)
	g.GET("/tweets/follow", h.GetTweets)
	g.GET("/users/:display_id/tweets", h.GetTweets)
	g.POST("/tweets/:tweet_id/retweet", h.Retweet)
	g.POST("/tweets/:tweet_id/reply", h.Retweet)
	g.POST("/tweet", h.CreateTweet)
}

// type GetAllTweetsResponseUser struct {
// 	ID        int    `json:"id"`
// 	DisplayID string `json:"display_id"`
// 	Name      string `json:"name"`
// }

// type GetAllTweetsResponseRetweet struct {
// 	ID        int                      `json:"id"`
// 	User      GetAllTweetsResponseUser `json:"user"`
// 	Content   string                   `json:"content"`
// 	CreatedAt time.Time                `json:"created_at"`
// }

// type GetAllTweetsResponseReply = GetAllTweetsResponseRetweet

// type GetAllTweetsResponse struct {
// 	ID        int                          `json:"id"`
// 	User      GetAllTweetsResponseUser     `json:"user"`
// 	Content   string                       `json:"content"`
// 	Retweet   *GetAllTweetsResponseRetweet `json:"retweet"`
// 	Reply     *GetAllTweetsResponseReply   `json:"reply"`
// 	CreatedAt time.Time                    `json:"created_at"`
// }

type GetTimelineTweetsResponseUser struct {
	ID        int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name      string `json:"name"`
	IconURL   string `json:"icon_url"`
}

type GetTimelineTweetsResponseRetweet struct {
	ID        int                           `json:"id"`
	User      GetTimelineTweetsResponseUser `json:"user"`
	Content   string                        `json:"content"`
	CreatedAt time.Time                     `json:"created_at"`
}

type GetTimelineTweetsResponseInteractions struct {
	LikeCount    int `json:"like_count"`
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
}

type GetTimelineTweetsResponse struct {
	ID           int                                   `json:"id"`
	User         GetTimelineTweetsResponseUser         `json:"user"`
	Content      string                                `json:"content"`
	Retweet      *GetTimelineTweetsResponseRetweet     `json:"retweet"`
	Interactions GetTimelineTweetsResponseInteractions `json:"interactions"`
	CreatedAt    time.Time                             `json:"created_at"`
}

// TODO: cursorを使ってページネーションする
func (h *TweetHandler) GetTimelineTweets(c echo.Context) error {
	var tweets []model.Tweet
	// TweetとUserをJOINして取得
	err := h.db.Select(&tweets, `
		SELECT tweets.*, user_profile.user_id as "user.id", user_profile.name as "user.name", user_profile.display_id as "user.display_id", user_profile.icon_url as "user.icon_url"
		FROM tweets
		JOIN user_profile ON tweets.user_id = user_profile.user_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	var retweets []model.Tweet
	var retweetIDs []int
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetIDs = append(retweetIDs, *tweet.RetweetID)
		}
	}
	// RetweetがあればRetweet対象のツイートを取得
	if len(retweetIDs) > 0 {
		query, args, err := sqlx.In(`
			SELECT tweets.*, user_profile.user_id as "user.id", user_profile.name as "user.name", user_profile.display_id as "user.display_id"
			FROM tweets
			JOIN user_profile ON tweets.user_id = user_profile.user_id
			WHERE tweets.id IN (?)
		`, retweetIDs)

		if err != nil {
			return h.handleError(c, err)
		}
		err = h.db.Select(&retweets, query, args...)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	// RetweetとRetweet対象のツイートを紐づける
	var retweetMap = map[int]model.Tweet{}
	for _, retweet := range retweets {
		retweetMap[retweet.ID] = retweet
	}
	for i, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweet, ok := retweetMap[*tweet.RetweetID]
			if !ok {
				return h.handleError(c, errors.New("retweet not found"))
			}
			tweets[i].Retweet = &retweet
		}
	}

	// いいね、リツイート、リプライの数を取得
	// それぞれ初期化
	likeCounts := map[int]int
	retweetCounts := map[int]int
	replyCounts := map[int]int

	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetCounts[*tweet.RetweetID]++
		}
		if tweet.ReplyID != nil {
			replyCounts[*tweet.ReplyID]++
		}
	}

	
	// レスポンス用の構造体に変換
	res := make([]GetTimelineTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetTimelineTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetTimelineTweetsResponseRetweet{
				ID: tweet.Retweet.ID,
				User: GetTimelineTweetsResponseUser{
					ID:        tweet.Retweet.User.ID,
					Name:      tweet.Retweet.User.Name,
					DisplayID: tweet.Retweet.User.DisplayID,
					IconURL:   tweet.Retweet.User.IconURL,
				},
				Content:   tweet.Retweet.Content,
				CreatedAt: tweet.Retweet.CreatedAt,
			}
		}
		res[i] = GetTimelineTweetsResponse{
			ID:        tweet.ID,
			User: GetTimelineTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
				IconURL:   tweet.User.IconURL,
			},
			Content:   tweet.Content,
			Retweet:   retweet,
			CreatedAt: tweet.CreatedAt,
		}
	}

func (h *TweetHandler) handleError(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(500, map[string]string{"message": "Internal Server Error"})
}

type GetTweetsResponseUser struct {
	ID        int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name      string `json:"name"`
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
	err := h.db.Select(&tweets, `
		SELECT tweets.*, user_profile.user_id as "user.id", user_profile.name as "user.name", user_profile.display_id as "user.display_id"
		FROM tweets
		JOIN user_profile ON tweets.user_id = user_profile.user_id
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
			SELECT tweets.*, user_profile.user_id as "user.id", user_profile.name as "user.name", user_profile.display_id as "user.display_id"
			FROM tweets
			JOIN user_profile ON tweets.user_id = user_profile.user_id
			WHERE tweets.id IN (?)
		`, retweetIDs)

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

	res := make([]GetTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetTweetsResponseRetweet{
				ID: tweet.Retweet.ID,
				User: GetTweetsResponseUser{
					ID:        tweet.Retweet.User.ID,
					Name:      tweet.Retweet.User.Name,
					DisplayID: tweet.Retweet.User.DisplayID,
				},
				Content: tweet.Retweet.Content,
			}
		}
		res[i] = GetTweetsResponse{
			ID: tweet.ID,
			User: GetTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
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
