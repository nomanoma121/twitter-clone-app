package handler

import (
	"database/sql"
	"errors"
	"log"
	"server/model"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{db: db, validator: validator.New()}
}

func (h *UserHandler) Register(g *echo.Group) {
	g.POST("/users/follow", h.Follow)
	g.GET("/users/followers", h.GetFollowers)
	// g.GET("/users/followees", h.GetFollowees)
}

type FollowRequest struct {
	UserID int `json:"user_id" validate:"required"`
}

func (h *UserHandler) Follow(c echo.Context) error {
	followerUserID := c.Get("user_id").(int)

	req := new(FollowRequest)
	if err := c.Bind(req); err != nil {
		log.Println(err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}
	if err := h.validator.Struct(req); err != nil {
		log.Println(err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	followeeUserID := req.UserID

	_, err := h.db.Exec("INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)", followerUserID, followeeUserID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(200)
}

type GetFollowersResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) GetFollowers(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var followers []model.User
	err := h.db.Select(&followers, `
		SELECT users.*
		FROM users
		JOIN follows ON users.id = follows.follower_id
		WHERE follows.followee_id = ?
	`, userID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetFollowersResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := make([]GetFollowersResponse, len(followers))
	for i, follower := range followers {
		res[i] = GetFollowersResponse{
			ID:    follower.ID,
			Name:  follower.Name,
			Email: follower.Email,
		}
	}

	return c.JSON(200, res)
}
