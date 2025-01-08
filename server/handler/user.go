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
	g.GET("/users/followees", h.GetFollowees)
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
	DisplayID string `json:"display_id"`
}

func (h *UserHandler) GetFollowers(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var followers []model.UserProfile
	err := h.db.Select(&followers, `
		SELECT user_profiles.user_id AS id, user_profiles.name, user_profiles.display_id
		FROM user_profiles
		JOIN follows ON user_profiles.user_id = follows.follower_id
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
			DisplayID: follower.DisplayID,
		}
	}

	return c.JSON(200, res)
}

type GetFolloweesResponse struct {
	ID    int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name  string `json:"name"`
}

func (h *UserHandler) GetFollowees(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var followees []model.UserProfile
	err := h.db.Select(&followees, `
		SELECT user_profiles.user_id AS id, user_profiles.name, user_profiles.display_id
		FROM user_profiles
		JOIN follows ON user_profiles.user_id = follows.followee_id
		WHERE follows.follower_id = ?
	`, userID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetFolloweesResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := make([]GetFolloweesResponse, len(followees))
	for i, followee := range followees {
		res[i] = GetFolloweesResponse{
			ID:    followee.ID,
			Name:  followee.Name,
			DisplayID: followee.DisplayID,
		}
	}

	return c.JSON(200, res)
}
