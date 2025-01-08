package handler

import (
	"database/sql"
	"errors"
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
	ID        int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name      string `json:"name"`
}

type GetAllTweetsResponseRetweet struct {
	ID        int                      `json:"id"`
	User      GetAllTweetsResponseUser `json:"user"`
	Content   string                   `json:"content"`
	CreatedAt time.Time                `json:"created_at"`
}

type GetAllTweetsResponseReply = GetAllTweetsResponseRetweet

type GetAllTweetsResponse struct {
	ID        int                          `json:"id"`
	User      GetAllTweetsResponseUser     `json:"user"`
	Content   string                       `json:"content"`
	Retweet   *GetAllTweetsResponseRetweet `json:"retweet"`
	Reply     *GetAllTweetsResponseReply   `json:"reply"`
	CreatedAt time.Time                    `json:"created_at"`
}

func (h *TweetHandler) GetAllTweets(c echo.Context) error {
	tweets, err := h.fetchTweets()
	if err != nil {
		return h.handleError(c, err)
	}

	retweetMap, replyMap, err := h.fetchRelatedTweets(tweets)
	if err != nil {
		return h.handleError(c, err)
	}

	h.populateRelatedTweets(tweets, retweetMap, replyMap)

	response := h.buildResponse(tweets)
	return c.JSON(200, response)
}

func (h *TweetHandler) fetchTweets() ([]model.Tweet, error) {
	var tweets []model.Tweet
	err := h.db.Select(&tweets, `
		SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.display_id as "user.display_id"
		FROM tweets
		JOIN users ON tweets.user_id = users.id
	`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return tweets, err
}

func (h *TweetHandler) fetchRelatedTweets(tweets []model.Tweet) (map[int]model.Tweet, map[int]model.Tweet, error) {
	retweetIDs, replyIDs := h.extractRelatedIDs(tweets)

	retweets, err := h.fetchTweetsByIDs(retweetIDs)
	if err != nil {
		return nil, nil, err
	}

	replies, err := h.fetchTweetsByIDs(replyIDs)
	if err != nil {
		return nil, nil, err
	}

	return h.buildTweetMap(retweets), h.buildTweetMap(replies), nil
}

func (h *TweetHandler) extractRelatedIDs(tweets []model.Tweet) ([]int, []int) {
	var retweetIDs, replyIDs []int
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetIDs = append(retweetIDs, *tweet.RetweetID)
		}
		if tweet.ReplyID != nil {
			replyIDs = append(replyIDs, *tweet.ReplyID)
		}
	}
	return retweetIDs, replyIDs
}

func (h *TweetHandler) fetchTweetsByIDs(ids []int) ([]model.Tweet, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`
		SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.display_id as "user.display_id"
		FROM tweets
		JOIN users ON tweets.user_id = users.id
		WHERE tweets.id IN (?)
	`, ids)
	if err != nil {
		return nil, err
	}

	var tweets []model.Tweet
	err = h.db.Select(&tweets, query, args...)
	return tweets, err
}

func (h *TweetHandler) buildTweetMap(tweets []model.Tweet) map[int]model.Tweet {
	result := make(map[int]model.Tweet)
	for _, tweet := range tweets {
		result[tweet.ID] = tweet
	}
	return result
}

func (h *TweetHandler) populateRelatedTweets(tweets []model.Tweet, retweetMap, replyMap map[int]model.Tweet) {
	for i, tweet := range tweets {
		if tweet.RetweetID != nil {
			tweets[i].Retweet = h.lookupTweet(retweetMap, *tweet.RetweetID)
		}
		if tweet.ReplyID != nil {
			tweets[i].Reply = h.lookupTweet(replyMap, *tweet.ReplyID)
		}
	}
}

func (h *TweetHandler) lookupTweet(tweetMap map[int]model.Tweet, id int) *model.Tweet {
	if tweet, exists := tweetMap[id]; exists {
		return &tweet
	}
	return nil
}

func (h *TweetHandler) buildResponse(tweets []model.Tweet) []GetAllTweetsResponse {
	response := make([]GetAllTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		response[i] = GetAllTweetsResponse{
			ID: tweet.ID,
			User: GetAllTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
			},
			Content: tweet.Content,
			Retweet: h.buildRetweetResponse(tweet.Retweet),
			Reply:   h.buildRetweetResponse(tweet.Reply),
			CreatedAt: tweet.CreatedAt,
		}
	}
	return response
}

func (h *TweetHandler) buildRetweetResponse(tweet *model.Tweet) *GetAllTweetsResponseRetweet {
	if tweet == nil {
		return nil
	}
	return &GetAllTweetsResponseRetweet{
		ID: tweet.ID,
		User: GetAllTweetsResponseUser{
			ID:        tweet.User.ID,
			Name:      tweet.User.Name,
			DisplayID: tweet.User.DisplayID,
		},
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt,
	}
}

func (h *TweetHandler) handleError(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(500, map[string]string{"message": "Internal Server Error"})
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

// Retweet or Replyがある場合はそれを取得し、Tweetに加える関数
func getRetweet(retweets *[]model.Tweet, tweets *[]model.Tweet) {
	var retweetsIDs []int
	for _, tweet := range *tweets {
		if tweet.RetweetID != nil {
			retweetsIDs = append(retweetsIDs, *tweet.RetweetID)
		}
	}
	if len(retweetsIDs) > 0 {
		// getTweetsWithUserでRetweetを取得

	var retweetMap = map[int]model.Tweet{}
	for _, retweet := range *retweets {
		retweetMap[retweet.ID] = retweet
	}
	for i, tweet := range *tweets {
		if tweet.RetweetID != nil {
			retweet, ok := retweetMap[*tweet.RetweetID]
			if !ok {
				return c.JSON(500, map[string]string{"message": "Internal Server Error"})
			}
			(*tweets)[i].Retweet = &retweet
		}
	}
}

// func getReply(r *model.Tweet, tweets *[]model.Tweet) {
// 	var replyIDs []int
// 	for _, tweet := range *tweets {
// 		if tweet.ReplyID != nil {
// 			replyIDs = append(replyIDs, *tweet.ReplyID)
// 		}
// 	}

// }

func (h *TweetHandler) getTweetsWithUser(tweetIDs []int, tweets *[]model.Tweet) error {
	query, args, err := sqlx.In(`
		SELECT tweets.*, users.id as "user.id", users.name as "user.name", users.email as "user.email"
		FROM tweets
		JOIN users ON tweets.user_id = users.id
		WHERE tweets.id IN (?)
	`, tweetIDs)
	if err != nil {
		log.Println(err)
		return err
	}
	err = h.db.Select(&tweets, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
