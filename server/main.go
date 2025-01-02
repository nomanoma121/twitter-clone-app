package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

const (
	APP_PORT = 8000
)

func main() {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")

	if user == "" || password == "" || dbname == "" || host == "" || port == "" {
		log.Fatal("環境変数が設定されていません")
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	// DB接続、5秒ごとに10回リトライ (MySQLコンテナが立ち上がるまで待つ)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Printf("DB接続に失敗しました: %v", err)
		for i := 0; i < 10; i++ {
			err = db.Ping()
			if err == nil {
				break
			}
			log.Printf("DB接続に失敗しました (%d): %v", i+1, err)
			time.Sleep(5 * time.Second)
		}
	}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	if mode == "development" {
		e.Debug = true
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", APP_PORT)))
}
