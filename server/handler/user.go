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
	g.POST("/users/:displayID/follow", h.Follow)
	g.DELETE("/users/:displayID/unfollow", h.Unfollow)
	g.GET("/users/:displayID/followers", h.GetFollowers)
	g.GET("/users/:displayID/followees", h.GetFollowees)
	g.GET("/users/:displayID", h.GetUser)
	g.GET("/users/:displayID/tweet-counts", h.GetTweetCounts)
}

type GetUserResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	DisplayID      string    `json:"display_id"`
	IconURL        string    `json:"icon_url"`
	HeaderURL      string    `json:"header_url"`
	Profile        string    `json:"profile"`
	FollowerCounts int       `json:"follower_counts"`
	FolloweeCounts int       `json:"followee_counts"`
	FollowedByUser bool      `json:"followed_by_user"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserData struct {
	ID             int           `db:"user_id"`
	Name           string        `db:"name"`
	DisplayID      string        `db:"display_id"`
	IconURL        string        `db:"icon_url"`
	HeaderURL      string        `db:"header_url"`
	Profile        string        `db:"profile"`
	FollowerCounts sql.NullInt64 `db:"follower_counts"`
	FolloweeCounts sql.NullInt64 `db:"followee_counts"`
	CreatedAt      time.Time     `db:"created_at"`
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID := c.Get("user_id").(int)
	displayID := c.Param("displayID")

	targetUserID, err := h.getUserIDByDisplayID(displayID)
	if err != nil {
		return c.JSON(404, map[string]string{"message": "User not found"})
	}

	var user UserData
	err = h.db.Get(&user, `
		SELECT user_profiles.user_id, user_profiles.name, user_profiles.display_id, user_profiles.icon_url, user_profiles.header_url, user_profiles.profile, user_profiles.created_at,
		(SELECT COUNT(*) FROM follows WHERE followee_id = user_profiles.user_id) AS followee_counts,
		(SELECT COUNT(*) FROM follows WHERE follower_id = user_profiles.user_id) AS follower_counts
		FROM user_profiles
		WHERE user_profiles.user_id = ?
	`, targetUserID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(404, map[string]string{"message": "User not found"})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	// あるユーザーがこのユーザーをフォローしているか
	isFollowed, err := h.isFollowed(userID, targetUserID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := GetUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		DisplayID:      user.DisplayID,
		IconURL:        user.IconURL,
		HeaderURL:      user.HeaderURL,
		Profile:        user.Profile,
		FollowerCounts: 0,
		FolloweeCounts: 0,
		FollowedByUser: isFollowed,
		CreatedAt:      user.CreatedAt,
	}

	if user.FollowerCounts.Valid {
		res.FollowerCounts = int(user.FollowerCounts.Int64)
	}

	if user.FolloweeCounts.Valid {
		res.FolloweeCounts = int(user.FolloweeCounts.Int64)
	}

	return c.JSON(200, res)
}

func (h *UserHandler) Follow(c echo.Context) error {
	followerUserID := c.Get("user_id").(int)
	displayID := c.Param("displayID")

	followeeUserID, err := h.getUserIDByDisplayID(displayID)
	if err != nil {
		return c.JSON(404, map[string]string{"message": "User not found"})
	}

	_, err = h.db.Exec("INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)", followerUserID, followeeUserID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(200)
}

func (h *UserHandler) Unfollow(c echo.Context) error {
	followerUserID := c.Get("user_id").(int)
	displayID := c.Param("displayID")

	followeeUserID, err := h.getUserIDByDisplayID(displayID)
	if err != nil {
		return c.JSON(404, map[string]string{"message": "User not found"})
	}


	_, err = h.db.Exec("DELETE FROM follows WHERE follower_id = ? AND followee_id = ?", followerUserID, followeeUserID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(200)
}

type GetFollowersResponse struct {
	ID             int    `json:"id"`
	DisplayID      string `json:"display_id"`
	Name           string `json:"name"`
	IconURL        string `json:"icon_url"`
	Profile        string `json:"profile"`
	FollowedByUser bool   `json:"followed_by_user"`
}

func (h *UserHandler) GetFollowers(c echo.Context) error {
	clientUserID := c.Get("user_id").(int)
	displayID := c.Param("displayID")

	targetUserID, err := h.getUserIDByDisplayID(displayID)

	var followers []model.UserProfile
	err = h.db.Select(&followers, `
		SELECT user_profiles.user_id AS id, user_profiles.name, user_profiles.display_id, user_profiles.icon_url, user_profiles.profile
		FROM user_profiles
		JOIN follows ON user_profiles.user_id = follows.follower_id
		WHERE follows.followee_id = ?
	`, targetUserID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetFollowersResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	var isFollowedMap = make([]bool, len(followers))
	for i, follower := range followers {
		isFollowedMap[i], err = h.isFollowed(clientUserID, follower.ID)
	}

	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := make([]GetFollowersResponse, len(followers))
	for i, follower := range followers {
		res[i] = GetFollowersResponse{
			ID:             follower.ID,
			Name:           follower.Name,
			DisplayID:      follower.DisplayID,
			IconURL:        follower.IconURL,
			Profile:        follower.Profile,
			FollowedByUser: isFollowedMap[i],
		}
	}

	return c.JSON(200, res)
}

type GetFolloweesResponse struct {
	ID             int    `json:"id"`
	DisplayID      string `json:"display_id"`
	Name           string `json:"name"`
	IconURL        string `json:"icon_url"`
	Profile        string `json:"profile"`
	FollowedByUser bool   `json:"followed_by_user"`
}

func (h *UserHandler) GetFollowees(c echo.Context) error {
	clientUserID := c.Get("user_id").(int)
	displayID := c.Param("displayID")

	targetUserID, err := h.getUserIDByDisplayID(displayID)

	var followees []model.UserProfile
	err = h.db.Select(&followees, `
		SELECT user_profiles.user_id AS id, user_profiles.name, user_profiles.display_id, user_profiles.icon_url, user_profiles.profile
		FROM user_profiles
		JOIN follows ON user_profiles.user_id = follows.followee_id
		WHERE follows.follower_id = ?
	`, targetUserID)

	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetFolloweesResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	var isFollowedMap = make([]bool, len(followees))
	for i, followee := range followees {
		isFollowedMap[i], err = h.isFollowed(clientUserID, followee.ID)
	}

	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res := make([]GetFolloweesResponse, len(followees))
	for i, followee := range followees {
		res[i] = GetFolloweesResponse{
			ID:             followee.ID,
			Name:           followee.Name,
			DisplayID:      followee.DisplayID,
			IconURL:        followee.IconURL,
			Profile:        followee.Profile,
			FollowedByUser: isFollowedMap[i],
		}
	}

	return c.JSON(200, res)
}

// 1 -> 2のフォロー関係のboolを返す
func (h *UserHandler) isFollowed(followerID, followeeID int) (bool, error) {
	var count int
	err := h.db.Get(&count, "SELECT COUNT(*) FROM follows WHERE follower_id = ? AND followee_id = ?", followerID, followeeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// display_id -> user_id
func (h *UserHandler) getUserIDByDisplayID(displayID string) (int, error) {
	var userID int
	err := h.db.Get(&userID, "SELECT user_id FROM user_profiles WHERE display_id = ?", displayID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

type GetTweetCountsResponse struct {
	TweetCounts int `json:"tweet_counts"`
}

func (h *UserHandler) GetTweetCounts(c echo.Context) error {
	displayID := c.Param("displayID")

	userID, err := h.getUserIDByDisplayID(displayID)
	if err != nil {
		return c.JSON(404, map[string]string{"message": "User not found"})
	}

	var tweetCounts int
	err = h.db.Get(&tweetCounts, "SELECT COUNT(*) FROM tweets WHERE user_id = ?", userID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, GetTweetCountsResponse{TweetCounts: tweetCounts})
}
