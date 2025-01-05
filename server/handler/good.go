package handler

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type GoodHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
}

func NewGoodHandler(db *sqlx.DB) *GoodHandler {
	return &GoodHandler{db: db, validator: validator.New()}
}

func (h *GoodHandler) Register(g *echo.Group) {
	g.POST("/good", h.Good)
	g.GET("/goods", h.GetGoods)
}

type GoodRequest struct {
	TweetID int `json:"tweet_id" validate:"required"`
}

func (h *GoodHandler) Good(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(GoodRequest)
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
	err := h.db.Get(&count, "SELECT COUNT(*) FROM goods WHERE tweet_id = ? AND user_id = ?", tweetID, userID)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	if count != 0 {
		return c.JSON(400, map[string]string{"message": "Already Good"})
	}

	_, err = h.db.Exec("INSERT INTO goods (tweet_id, user_id) VALUES (?, ?)", tweetID, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type GetGoodsResponse struct {
	ID    int `json:"id"`
	Tweet struct {
		ID   int `json:"id"`
		User struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"user"`
		Content *string `json:"content"`
	} `json:"tweet"`
}

func (h *GoodHandler) GetGoods(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var goods []GetGoodsResponse
	err := h.db.Select(&goods, `
		SELECT 
    goods.id AS "id",
    tweets.id AS "tweet.id",
    tweets.content AS "tweet.content",
    users.id AS "tweet.user.id",
    users.name AS "tweet.user.name",
    users.email AS "tweet.user.email"
		FROM goods
		JOIN tweets ON goods.tweet_id = tweets.id
		JOIN users ON tweets.user_id = users.id
		WHERE goods.user_id = ?
	`, userID)

	log.Printf("goods: %#v\n", goods)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, goods)
}
