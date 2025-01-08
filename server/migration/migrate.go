package migration

import (
	"github.com/jmoiron/sqlx"
)

// users
func Migrate(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// user_profilesテーブルを作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_profiles (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			display_id VARCHAR(255) NOT NULL UNIQUE,
			icon_url VARCHAR(255),
			header_url VARCHAR(255),
			profile TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tweets (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			content TEXT NOT NULL,
			retweet_id INT,
			reply_id INT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// tweet_imagesテーブルを作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tweet_images (
			id INT PRIMARY KEY AUTO_INCREMENT,
			tweet_id INT NOT NULL,
			url VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS follows (
			id INT PRIMARY KEY AUTO_INCREMENT,
			follower_id INT NOT NULL,
			followee_id INT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE(follower_id, followee_id)
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			tweet_id INT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

var tables = []string{
	"users",
	"user_profiles",
	"tweets",
	"tweet_images",
	"follows",
	"likes",
	"todos",
}

func Truncate(db *sqlx.DB) error {
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			return err
		}
	}

	return nil
}

func Reset(db *sqlx.DB) error {
	if err := Truncate(db); err != nil {
		return err
	}

	if err := Migrate(db); err != nil {
		return err
	}

	return nil
}
