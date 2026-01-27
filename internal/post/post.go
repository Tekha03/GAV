package post

type Post struct {
	ID		uint	`json:"id"`
	UserID 	uint	`json:"user_id"`
	Content string	`json:"content"`
}
