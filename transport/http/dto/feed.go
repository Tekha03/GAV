package dto

type FeedResponse struct {
	Posts		[]PostResponse	`json:"posts"`
	NextCursor	string			`json:"next_cursor,omitempty"`
	HasMore		bool			`json:"has_more"`
}
