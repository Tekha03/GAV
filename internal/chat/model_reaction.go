package chat

type Reaction struct {
	ID        uint
	MessageID uint
	UserID    uint
	Emoji     string
}
