package profile

type UserProfile struct {
    UserID       uint   `gorm:"primaryKey"`
    Name         string
    Surname      string
    Username     string `gorm:"unique"`
    ProfilePhoto string
    Bio          string
    Address      string
    BirthDate    string
}
