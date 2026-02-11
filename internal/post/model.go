package post

import "time"

type Post struct {
	ID			uint		`json:"id"`
	UserID 		uint		`json:"user_id"`
	Content 	string		`json:"content"`
	CreatedAt	time.Time	`json:"created_at"`
}
