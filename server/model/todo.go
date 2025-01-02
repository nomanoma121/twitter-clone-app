package model

import "time"

type Todo struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Title     string    `db:"title"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
