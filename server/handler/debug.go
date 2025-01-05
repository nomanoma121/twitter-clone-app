package handler

import (
	"log"
	"server/migration"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type DebugHandler struct {
	db *sqlx.DB
}

func NewDebugHandler(db *sqlx.DB) *DebugHandler {
	return &DebugHandler{db: db}
}

func (h *DebugHandler) Register(g *echo.Group) {
	g.POST("/seed", h.Seed)
}

func (h *DebugHandler) Seed(c echo.Context) error {
	err := migration.Seed(h.db)
	if err != nil {
		log.Println(err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.NoContent(200)
}
