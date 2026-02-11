package chat

import "time"

type Chat struct {
	ID			uint
	IsGroup		bool
	Title		string
	CreatedAt 	time.Time
}