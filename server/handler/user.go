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
	g.GET("/users/:id", h.GetUser)
}

type GetUserResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DisplayID      string `json:"display_id"`
	IconURL        string `json:"icon_url"`
	HeaderURL      string `json:"header_url"`
	Profile        string `json:"profile"`
	FollowerCounts int    `json:"follower_counts"`
	FolloweeCounts int    `json:"followee_counts"`
	CreatedAt      string `json:"created_at"`
}

type UserData struct {
	ID             int       `db:"user_id"`
	Name           string    `db:"name"`
	DisplayID      string    `db:"display_id"`
	IconURL        string    `db:"icon_url"`
	HeaderURL      string    `db:"header_url"`
	Profile        string    `db:"profile"`
	FollowerCounts int       `db:"follower_counts"`
	FolloweeCounts int       `db:"followee_counts"`
	CreatedAt      time.Time `db:"created_at"`
}

func (h *UserHandler) GetUser(c echo.Context) error {
	display_id := c.Param("id")

	var user UserData
	// TODO: なんか冗長な気がするからあとでリファクタリングする
	err := h.db.Get(&user, `
		SELECT user_profiles.user_id, user_profiles.name, user_profiles.display_id, user_profiles.icon_url, user_profiles.header_url, user_profiles.profile, follower_counts, followee_counts, user_profiles.created_at
		FROM user_profiles
		LEFT JOIN (
			SELECT followee_id, COUNT(follower_id) AS follower_counts
			FROM follows
			GROUP BY followee_id
		) AS followers ON user_profiles.user_id = followers.followee_id
		LEFT JOIN (
			SELECT follower_id, COUNT(followee_id) AS followee_counts
			FROM follows
			GROUP BY follower_id
		) AS followees ON user_profiles.user_id = followees.follower_id
		WHERE user_profiles.display_id = ?
	`, display_id)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "User not found"})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := GetUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		DisplayID:      user.DisplayID,
		IconURL:        user.IconURL,
		HeaderURL:      user.HeaderURL,
		Profile:        user.Profile,
		FollowerCounts: user.FollowerCounts,
		FolloweeCounts: user.FolloweeCounts,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(200, res)
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
	ID        int    `json:"id"`
	Name      string `json:"name"`
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
			ID:        follower.ID,
			Name:      follower.Name,
			DisplayID: follower.DisplayID,
		}
	}

	return c.JSON(200, res)
}

type GetFolloweesResponse struct {
	ID        int    `json:"id"`
	DisplayID string `json:"display_id"`
	Name      string `json:"name"`
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
			ID:        followee.ID,
			Name:      followee.Name,
			DisplayID: followee.DisplayID,
		}
	}

	return c.JSON(200, res)
}
