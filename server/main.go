package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	APP_PORT = 8000
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", APP_PORT)))
}
