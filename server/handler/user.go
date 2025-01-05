package handler

import (
	"log"

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
