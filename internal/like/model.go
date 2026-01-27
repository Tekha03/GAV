package like

type Like struct {
	UserID	uint	`gorm:"primaryKey"`
	PostID	uint	`gorm:"primaryKey"`
}
