package main

import (
	"fmt"
	"log"
	"os"
	"server/handler"
	"server/middleware"
	"server/migration"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

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
		log.Fatal("DB接続情報の環境変数が設定されていません")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	// DB接続、1秒ごとに30回リトライ (MySQLコンテナが立ち上がるまで待つ)
	var db *sqlx.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sqlx.Connect("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("DB接続に失敗しました (%d): %v", i+1, err)
		time.Sleep(1 * time.Second)
	}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("環境変数[JWT_SECRET]が設定されていません")
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development"
	}

	origin := os.Getenv("ORIGIN")
	if mode == "production" {
		if origin == "" {
			log.Fatal("環境変数[ORIGIN]が設定されていません")
		}
	}
	e := echo.New()

	if mode == "development" {
		e.Debug = true
		e.Use(echoMiddleware.Logger())
		e.Use(echoMiddleware.Recover())
		e.Use(echoMiddleware.CORS())
	} else {
		e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{origin},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	authMiddleware := middleware.NewAuthMiddleware(db, jwtSecret)

	authGroup := e.Group("/auth")
	loginHandler := handler.NewAuthHandler(db, jwtSecret)
	loginHandler.Register(authGroup, authMiddleware)

	apiGroup := e.Group("/api")
	apiGroup.Use(authMiddleware.Middleware())
	todoHandler := handler.NewTodoHandler(db)
	todoHandler.Register(apiGroup)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", APP_PORT)))
}
