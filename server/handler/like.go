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
	g.POST("/like/:id", h.CreateLike)
	g.DELETE("/like/:id", h.DeleteLike)
}

func (h *LikeHandler) CreateLike(c echo.Context) error {
	userID := c.Get("user_id").(int)
	tweetID := c.Param("id")

	log.Printf("userID: %d, tweetID: %s", userID, tweetID)

	_, err := h.db.Exec("INSERT INTO likes (user_id, tweet_id) VALUES (?, ?)", userID, tweetID)
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}

	return c.NoContent(200)
}

func (h *LikeHandler) DeleteLike(c echo.Context) error {
	userID := c.Get("user_id").(int)
	tweetID := c.Param("id")

	_, err := h.db.Exec("DELETE FROM likes WHERE user_id = ? AND tweet_id = ?", userID, tweetID)
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}

	return c.NoContent(200)
}
