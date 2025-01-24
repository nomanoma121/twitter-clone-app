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
	g.GET("/tweets/timeline", h.GetTimelineTweets)
	g.GET("/tweets/follow", h.GetFollowTweets)
	g.GET("/tweet/:id", h.GetTweetByID)
	g.GET("/tweets/:id/replies", h.GetTweetReplies)
	g.GET("/users/:display_id/tweets", h.GetUserTweets)
	g.POST("/tweet", h.CreateTweet)
	g.POST("/tweet/:id/retweet", h.CreateRetweet)
	g.POST("/tweet/:id/reply", h.CreateReply)
}

type GetTimelineTweetsResponseUser struct {
	ID        int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name      string `json:"name"`
	IconURL   string `json:"icon_url"`
}

type GetTimelineTweetsResponseRetweet struct {
	ID            int                                   `json:"id"`
	User          GetTimelineTweetsResponseUser         `json:"user"`
	Content       string                                `json:"content"`
	Interactions  GetTimelineTweetsResponseInteractions `json:"interactions"`
	Liked_By_User bool                                  `json:"liked_by_user"`
	CreatedAt     time.Time                             `json:"created_at"`
}

type GetTimelineTweetsResponseInteractions struct {
	LikeCount    int `json:"like_count"`
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
}

// TODO: Userがいいねしているかどうかをboolで返す
type GetTimelineTweetsResponse struct {
	ID            int                                   `json:"id"`
	User          GetTimelineTweetsResponseUser         `json:"user"`
	Content       string                                `json:"content"`
	Retweet       *GetTimelineTweetsResponseRetweet     `json:"retweet"`
	Interactions  GetTimelineTweetsResponseInteractions `json:"interactions"`
	Liked_By_User bool                                  `json:"liked_by_user"`
	CreatedAt     time.Time                             `json:"created_at"`
}

// TODO: cursor, limitを使ってページネーションする
// HACK: コードが冗長なのでリファクタリングする
func (h *TweetHandler) GetTimelineTweets(c echo.Context) error {
	userID := c.Get("user_id").(int)
	var tweets []model.Tweet
	// TweetとUserをJOINして取得
	err := h.db.Select(&tweets, `
		SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
		FROM tweets
		JOIN user_profiles ON tweets.user_id = user_profiles.user_id
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
			SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
			FROM tweets
			JOIN user_profiles ON tweets.user_id = user_profiles.user_id
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
	var likeCounts []model.CountResult

	// いいね数を取得
	err = h.db.Select(&likeCounts, `
		SELECT tweet_id, COUNT(*) as count
		FROM likes
		GROUP BY tweet_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	var isLikedByUser = make([]bool, len(tweets))
	for i, tweet := range tweets {
		isLiked, err := h.isLiked(userID, tweet.ID)
		if err != nil {
			return h.handleError(c, err)
		}
		isLikedByUser[i] = isLiked
	}

	// マップに変換
	likeCountMap := map[int]int{}
	for _, count := range likeCounts {
		likeCountMap[count.TweetID] = count.Count
	}

	var retweetCountMap = map[int]int{}
	var replyCountMap = map[int]int{}

	for _, tweet := range tweets {
		log.Printf("reply_id: %v\n", tweet.ReplyID)
		if tweet.RetweetID != nil {
			retweetCountMap[*tweet.RetweetID]++
		}
		if tweet.ReplyID != nil {
			replyCountMap[*tweet.ReplyID]++
		}
	}

	log.Println(likeCountMap)
	log.Println(retweetCountMap)
	log.Println(replyCountMap)

	// tweetからreply_idのあるものを除外
	for i := 0; i < len(tweets); i++ {
		if tweets[i].ReplyID != nil {
			tweets = append(tweets[:i], tweets[i+1:]...)
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
				Content: tweet.Retweet.Content,
				Interactions: GetTimelineTweetsResponseInteractions{
					LikeCount:    likeCountMap[tweet.Retweet.ID],
					RetweetCount: retweetCountMap[tweet.Retweet.ID],
					ReplyCount:   replyCountMap[tweet.Retweet.ID],
				},
				CreatedAt: tweet.Retweet.CreatedAt,
			}
		}
		res[i] = GetTimelineTweetsResponse{
			ID: tweet.ID,
			User: GetTimelineTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
				IconURL:   tweet.User.IconURL,
			},
			Content: tweet.Content,
			Retweet: retweet,
			Interactions: GetTimelineTweetsResponseInteractions{
				LikeCount:    likeCountMap[tweet.ID],
				RetweetCount: retweetCountMap[tweet.ID],
				ReplyCount:   replyCountMap[tweet.ID],
			},
			Liked_By_User: isLikedByUser[i],
			CreatedAt:     tweet.CreatedAt,
		}
	}

	return c.JSON(200, res)
}

type GetFollowTweetsResponse = GetTimelineTweetsResponse
type GetFollowTweetsResponseRetweet = GetTimelineTweetsResponseRetweet
type GetFollowTweetsResponseInteractions = GetTimelineTweetsResponseInteractions
type GetFollowTweetsResponseUser = GetTimelineTweetsResponseUser

func (h *TweetHandler) GetFollowTweets(c echo.Context) error {
	userID := c.Get("user_id").(int)

	// フォローしているユーザーIDを取得
	var followIDs []int
	err := h.db.Select(&followIDs, "SELECT followee_id FROM follows WHERE follower_id = ?", userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// フォローしているユーザーのツイートを取得
	var tweets []model.Tweet
	if len(followIDs) > 0 {
		// replyは除外
		query, args, err := sqlx.In(`
			SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
			FROM tweets
			JOIN user_profiles ON tweets.user_id = user_profiles.user_id
			WHERE tweets.user_id IN (?)
			AND tweets.reply_id IS NULL
		`, followIDs)
		if err != nil {
			return h.handleError(c, err)
		}

		query = h.db.Rebind(query)
		err = h.db.Select(&tweets, query, args...)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	// リツイート情報の取得
	var retweets []model.Tweet
	var retweetIDs []int // intだとエラーになる
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetIDs = append(retweetIDs, *tweet.RetweetID)
		}
	}

	if len(retweetIDs) > 0 {
		query, args, err := sqlx.In(`
			SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
			FROM tweets
			JOIN user_profiles ON tweets.user_id = user_profiles.user_id
			WHERE tweets.id IN (?)
		`, retweetIDs)
		if err != nil {
			return h.handleError(c, err)
		}

		query = h.db.Rebind(query)
		err = h.db.Select(&retweets, query, args...)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	// リツイートIDをマップ化
	var retweetMap = map[int]model.Tweet{}
	for _, retweet := range retweets {
		retweetMap[retweet.ID] = retweet
	}

	// ツイートにリツイート情報を紐付け
	for i, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweet, ok := retweetMap[*tweet.RetweetID]
			if !ok {
				return h.handleError(c, errors.New("retweet not found"))
			}
			tweets[i].Retweet = &retweet
		}
	}

	// いいね数を取得
	var likeCounts []model.CountResult
	err = h.db.Select(&likeCounts, `
		SELECT tweet_id, COUNT(*) as count
		FROM likes
		GROUP BY tweet_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	likeCountMap := map[int]int{}
	for _, count := range likeCounts {
		likeCountMap[count.TweetID] = count.Count
	}

	// リツイート数と返信数をカウント
	var retweetCountMap = map[int]int{}
	var replyCountMap = map[int]int{}
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetCountMap[*tweet.RetweetID]++
		}
		if tweet.ReplyID != nil {
			replyCountMap[*tweet.ReplyID]++
		}
	}

	// いいねしているかどうかを取得
	var isLikedByUser = make([]bool, len(tweets))
	for i, tweet := range tweets {
		isLiked, err := h.isLiked(userID, tweet.ID)
		if err != nil {
			return h.handleError(c, err)
		}
		isLikedByUser[i] = isLiked
	}

	// レスポンスデータの作成
	res := make([]GetFollowTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetFollowTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetFollowTweetsResponseRetweet{
				ID: tweet.Retweet.ID,
				User: GetFollowTweetsResponseUser{
					ID:        tweet.Retweet.User.ID,
					Name:      tweet.Retweet.User.Name,
					DisplayID: tweet.Retweet.User.DisplayID,
					IconURL:   tweet.Retweet.User.IconURL,
				},
				Content: tweet.Retweet.Content,
				Interactions: GetFollowTweetsResponseInteractions{
					LikeCount:    likeCountMap[tweet.Retweet.ID],
					RetweetCount: retweetCountMap[tweet.Retweet.ID],
					ReplyCount:   replyCountMap[tweet.Retweet.ID],
				},
				CreatedAt: tweet.Retweet.CreatedAt,
			}
		}
		res[i] = GetFollowTweetsResponse{
			ID: tweet.ID,
			User: GetFollowTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
				IconURL:   tweet.User.IconURL,
			},
			Content: tweet.Content,
			Retweet: retweet,
			Interactions: GetFollowTweetsResponseInteractions{
				LikeCount:    likeCountMap[tweet.ID],
				RetweetCount: retweetCountMap[tweet.ID],
				ReplyCount:   replyCountMap[tweet.ID],
			},
			Liked_By_User: isLikedByUser[i],
			CreatedAt:     tweet.CreatedAt,
		}
	}

	return c.JSON(200, res)
}

type GetTweetByIDResponseUser = GetTimelineTweetsResponseUser
type GetTweetByIDResponseInteractions = GetTimelineTweetsResponseInteractions
type GetTweetByIDResponseRetweet = GetTimelineTweetsResponseRetweet
type GetTweetByIDResponse struct {
	ID            int                              `json:"id"`
	User          GetTweetByIDResponseUser         `json:"user"`
	Content       string                           `json:"content"`
	Retweet       *GetTweetByIDResponseRetweet     `json:"retweet"`
	Interactions  GetTweetByIDResponseInteractions `json:"interactions"`
	IsLikedByUser bool                             `json:"liked_by_user"`
	CreatedAt     time.Time                        `json:"created_at"`
}

func (h *TweetHandler) GetTweetByID(c echo.Context) error {
	userID := c.Get("user_id").(int)
	tweetID := c.Param("id")

	var tweet model.Tweet
	err := h.db.Get(&tweet, `
		SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
		FROM tweets
		JOIN user_profiles ON tweets.user_id = user_profiles.user_id
		WHERE tweets.id = ?
	`, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "Not Found"})
		}
		return h.handleError(c, err)
	}

	var retweet model.Tweet // ポインタから値にしたらエラーが消えた
	if tweet.RetweetID != nil {
		err = h.db.Get(&retweet, `
			SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
			FROM tweets
			JOIN user_profiles ON tweets.user_id = user_profiles.user_id
			WHERE tweets.id = ?
		`, *tweet.RetweetID)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	// いいね数を取得
	var likeCounts []model.CountResult
	err = h.db.Select(&likeCounts, `
		SELECT tweet_id, COUNT(*) as count
		FROM likes
		GROUP BY tweet_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	likeCountMap := map[int]int{}
	for _, count := range likeCounts {
		likeCountMap[count.TweetID] = count.Count
	}

	// リツイート数と返信数をカウント
	var retweetCountMap = map[int]int{}
	var replyCountMap = map[int]int{}
	if tweet.RetweetID != nil {
		retweetCountMap[*tweet.RetweetID]++
	}
	if tweet.ReplyID != nil {
		replyCountMap[*tweet.ReplyID]++
	}

	// いいねしいいねしているかどうかを取得
	isLiked, err := h.isLiked(userID, tweet.ID)
	if err != nil {
		return h.handleError(c, err)
	}

	res := GetTweetByIDResponse{
		ID: tweet.ID,
		User: GetTweetByIDResponseUser{
			ID:        tweet.User.ID,
			Name:      tweet.User.Name,
			DisplayID: tweet.User.DisplayID,
			IconURL:   tweet.User.IconURL,
		},
		Content: tweet.Content,
		Retweet: (*GetTweetByIDResponseRetweet)(nil),
		Interactions: GetTweetByIDResponseInteractions{
			LikeCount:    likeCountMap[tweet.ID],
			RetweetCount: retweetCountMap[tweet.ID],
			ReplyCount:   replyCountMap[tweet.ID],
		},
		IsLikedByUser: isLiked,
		CreatedAt:     tweet.CreatedAt,
	}

	if tweet.RetweetID != nil {
		res.Retweet = &GetTweetByIDResponseRetweet{
			ID: retweet.ID,
			User: GetTweetByIDResponseUser{
				ID:        retweet.User.ID,
				Name:      retweet.User.Name,
				DisplayID: retweet.User.DisplayID,
				IconURL:   retweet.User.IconURL,
			},
			Content: retweet.Content,
			Interactions: GetTweetByIDResponseInteractions{
				LikeCount:    likeCountMap[retweet.ID],
				RetweetCount: retweetCountMap[retweet.ID],
				ReplyCount:   replyCountMap[retweet.ID],
			},
			CreatedAt: retweet.CreatedAt,
		}
	}

	return c.JSON(200, res)
}

type GetTweetReplyResponseUser = GetTimelineTweetsResponseUser
type GetTweetReplyResponseInteractions = GetTimelineTweetsResponseInteractions

type GetTweetReplyResponse struct {
	ID            int                               `json:"id"`
	User          GetTweetReplyResponseUser         `json:"user"`
	Content       string                            `json:"content"`
	Interactions  GetTweetReplyResponseInteractions `json:"interactions"`
	IsLikedByUser bool                              `json:"liked_by_user"`
	CreatedAt     time.Time                         `json:"created_at"`
}

func (h *TweetHandler) GetTweetReplies(c echo.Context) error {
	tweetID := c.Param("id")
	userID := c.Get("user_id").(int)

	var tweets []model.Tweet
	err := h.db.Select(&tweets, `
		SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
		FROM tweets
		JOIN user_profiles ON tweets.user_id = user_profiles.user_id
		WHERE tweets.reply_id = ?
	`, tweetID)
	if err != nil {
		return h.handleError(c, err)
	}

	// いいね数を取得
	var likeCounts []model.CountResult
	err = h.db.Select(&likeCounts, `
		SELECT tweet_id, COUNT(*) as count
		FROM likes
		GROUP BY tweet_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	likeCountMap := map[int]int{}
	for _, count := range likeCounts {
		likeCountMap[count.TweetID] = count.Count
	}

	// リツイート数と返信数をカウント
	var retweetCountMap = map[int]int{}
	var replyCountMap = map[int]int{}
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetCountMap[*tweet.RetweetID]++
		}
		if tweet.ReplyID != nil {
			replyCountMap[*tweet.ReplyID]++
		}
	}

	isLikedByUser := make([]bool, len(tweets))
	for i, tweet := range tweets {
		isLiked, err := h.isLiked(userID, tweet.ID)
		if err != nil {
			return h.handleError(c, err)
		}
		isLikedByUser[i] = isLiked
	}

	res := make([]GetTweetReplyResponse, len(tweets))
	for i, tweet := range tweets {
		res[i] = GetTweetReplyResponse{
			ID: tweet.ID,
			User: GetTweetReplyResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
				IconURL:   tweet.User.IconURL,
			},
			Content: tweet.Content,
			Interactions: GetTweetReplyResponseInteractions{
				LikeCount:    likeCountMap[tweet.ID],
				RetweetCount: retweetCountMap[tweet.ID],
				ReplyCount:   replyCountMap[tweet.ID],
			},
			IsLikedByUser: isLikedByUser[i],
			CreatedAt:     tweet.CreatedAt,
		}
	}

	return c.JSON(200, res)
}

type GetUserTweetsResponseUser = GetTimelineTweetsResponseUser
type GetUserTweetsResponseInteractions = GetTimelineTweetsResponseInteractions
type GetUserTweetsResponseRetweet = GetTimelineTweetsResponseRetweet

type GetUserTweetsResponse struct {
	ID            int                               `json:"id"`
	User          GetUserTweetsResponseUser         `json:"user"`
	Content       string                            `json:"content"`
	Retweet       *GetUserTweetsResponseRetweet     `json:"retweet"`
	Interactions  GetUserTweetsResponseInteractions `json:"interactions"`
	IsLikedByUser bool                              `json:"is_liked_by_user"`
	CreatedAt     time.Time                         `json:"created_at"`
}

func (h *TweetHandler) GetUserTweets(c echo.Context) error {
	displayID := c.Param("display_id")
	userID := c.Get("user_id").(int)

	var user model.UserProfile
	err := h.db.Get(&user, "SELECT * FROM user_profiles WHERE display_id = ?", displayID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "Not Found"})
		}
		return h.handleError(c, err)
	}

	var tweets []model.Tweet
	err = h.db.Select(&tweets, `
		SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
		FROM tweets
		JOIN user_profiles ON tweets.user_id = user_profiles.user_id
		WHERE tweets.user_id = ?
	`, user.ID)
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

	if len(retweetIDs) > 0 {
		query, args, err := sqlx.In(`
			SELECT tweets.*, user_profiles.user_id as "user.id", user_profiles.name as "user.name", user_profiles.display_id as "user.display_id", user_profiles.icon_url as "user.icon_url"
			FROM tweets
			JOIN user_profiles ON tweets.user_id = user_profiles.user_id
			WHERE tweets.id IN (?)
		`, retweetIDs)
		if err != nil {
			return h.handleError(c, err)
		}

		query = h.db.Rebind(query)
		err = h.db.Select(&retweets, query, args...)
		if err != nil {
			return h.handleError(c, err)
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
				return h.handleError(c, errors.New("retweet not found"))
			}
			tweets[i].Retweet = &retweet
		}
	}

	var likeCounts []model.CountResult
	err = h.db.Select(&likeCounts, `
		SELECT tweet_id, COUNT(*) as count
		FROM likes
		GROUP BY tweet_id
	`)
	if err != nil {
		return h.handleError(c, err)
	}

	likeCountMap := map[int]int{}
	for _, count := range likeCounts {
		likeCountMap[count.TweetID] = count.Count
	}

	var retweetCountMap = map[int]int{}
	var replyCountMap = map[int]int{}
	for _, tweet := range tweets {
		if tweet.RetweetID != nil {
			retweetCountMap[*tweet.RetweetID]++
		}
		if tweet.ReplyID != nil {
			replyCountMap[*tweet.ReplyID]++
		}
	}

	isLikedByUser := make([]bool, len(tweets))
	for i, tweet := range tweets {
		isLiked, err := h.isLiked(userID, tweet.ID)
		if err != nil {
			return h.handleError(c, err)
		}
		isLikedByUser[i] = isLiked
	}

	res := make([]GetUserTweetsResponse, len(tweets))
	for i, tweet := range tweets {
		retweet := (*GetUserTweetsResponseRetweet)(nil)
		if tweet.Retweet != nil {
			retweet = &GetUserTweetsResponseRetweet{
				ID: tweet.Retweet.ID,
				User: GetUserTweetsResponseUser{
					ID:        tweet.Retweet.User.ID,
					Name:      tweet.Retweet.User.Name,
					DisplayID: tweet.Retweet.User.DisplayID,
					IconURL:   tweet.Retweet.User.IconURL,
				},
				Content: tweet.Retweet.Content,
				Interactions: GetUserTweetsResponseInteractions{
					LikeCount:    likeCountMap[tweet.Retweet.ID],
					RetweetCount: retweetCountMap[tweet.Retweet.ID],
					ReplyCount:   replyCountMap[tweet.Retweet.ID],
				},
				CreatedAt: tweet.Retweet.CreatedAt,
			}
		}
		res[i] = GetUserTweetsResponse{
			ID: tweet.ID,
			User: GetUserTweetsResponseUser{
				ID:        tweet.User.ID,
				Name:      tweet.User.Name,
				DisplayID: tweet.User.DisplayID,
				IconURL:   tweet.User.IconURL,
			},
			Content: tweet.Content,
			Retweet: retweet,
			Interactions: GetUserTweetsResponseInteractions{
				LikeCount:    likeCountMap[tweet.ID],
				RetweetCount: retweetCountMap[tweet.ID],
				ReplyCount:   replyCountMap[tweet.ID],
			},
			IsLikedByUser: isLikedByUser[i],
			CreatedAt:     tweet.CreatedAt,
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
	Content *string `json:"content"` // リツイートはcontentなくてもいい
}

func (h *TweetHandler) CreateRetweet(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(RetweetRequest)
	retweetID := c.Param("id")
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	log.Printf("req: %#v\n", *req.Content)
	log.Printf("retweetID: %v\n", retweetID)

	var tweet model.Tweet
	err := h.db.Get(&tweet, "SELECT * FROM tweets WHERE id = ?", retweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "Not Found"})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	if req.Content == nil {
		_, err = h.db.Exec("INSERT INTO tweets (user_id, retweet_id) VALUES (?, ?)", userID, retweetID)
		if err != nil {
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
	} else {
		_, err = h.db.Exec("INSERT INTO tweets (user_id, retweet_id, content) VALUES (?, ?, ?)", userID, retweetID, *req.Content)
		if err != nil {
			return c.JSON(500, map[string]string{"message": "Internal Server Error"})
		}
	}

	return c.NoContent(201)
}

type ReplyRequest struct {
	Content string `json:"content" validate:"required"`
}

func (h *TweetHandler) CreateReply(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(ReplyRequest)
	// パスからreply_idを取得
	replyID := c.Param("id")
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	log.Printf("req: %#v\n", req.Content)

	var tweet model.Tweet
	// reply_idが存在するか確認
	err := h.db.Get(&tweet, "SELECT * FROM tweets WHERE id = ?", replyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "Not Found"})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	_, err = h.db.Exec("INSERT INTO tweets (user_id, reply_id, content) VALUES (?, ?, ?)", userID, replyID, req.Content)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(201)
}

func (h *TweetHandler) handleError(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(500, map[string]string{"message": "Internal Server Error"})
}

// あるユーザーがあるツイートをいいねしているかどうかを返す
func (h *TweetHandler) isLiked(userID int, tweetID int) (bool, error) {
	var count int
	err := h.db.Get(&count, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND tweet_id = ?", userID, tweetID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
