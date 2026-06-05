package settings

type UpdateSettingsInput struct {
	ProfilePrivacy *bool `json:"profile_privacy"`
	ShowLocation   *bool `json:"show_location"`
	AllowMessages  *bool `json:"allow_messages"`
}
