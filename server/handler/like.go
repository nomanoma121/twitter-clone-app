package handler

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type LikeHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
}

func NewLikeHandler(db *sqlx.DB) *LikeHandler {
	return &LikeHandler{db: db, validator: validator.New()}
}

func (h *LikeHandler) Register(g *echo.Group) {
	g.POST("/like", h.Like)
	g.GET("/likes", h.GetLikes)
}

type LikeRequest struct {
	TweetID int `json:"tweet_id" validate:"required"`
}

func (h *LikeHandler) Like(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(LikeRequest)
	if err := c.Bind(req); err != nil {
		log.Println(err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}
	if err := h.validator.Struct(req); err != nil {
		log.Println(err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	tweetID := req.TweetID

	var count int
	err := h.db.Get(&count, "SELECT COUNT(*) FROM likes WHERE tweet_id = ? AND user_id = ?", tweetID, userID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	if count != 0 {
		return c.JSON(400, map[string]string{"message": "Already Good"})
	}

	_, err = h.db.Exec("INSERT INTO likes (tweet_id, user_id) VALUES (?, ?)", tweetID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type GetLikesResponse struct {
	ID    int `json:"id"`
	Tweet struct {
		ID   int `json:"id"`
		User struct {
			ID        int    `json:"id"`
			DisplayID string `json:"display_id"`
			Name      string `json:"name"`
		} `json:"user"`
		Content *string `json:"content"`
	} `json:"tweet"`
}

func (h *LikeHandler) GetLikes(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var likes []GetLikesResponse
	err := h.db.Select(&likes, `
		SELECT likes.id, tweets.*, users.display_id, users.name
		FROM likes	
		JOIN tweets ON likes.tweet_id = tweets.id
		JOIN users ON tweets.user_id = users.id
		WHERE likes.user_id = ?
	`, userID)

	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, likes)
}
