package settings

type Service interface {
    Get(userID uint) (*UserSettings, error)
    Update(userID uint, input UpdateSettingsInput) error
}
