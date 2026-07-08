package user

type LocationStatus uint8

const (
	Inactive LocationStatus = iota
	Walking
	ForcedOffline
)

type LocationVisibility uint8

const (
	VisibilityEveryone LocationVisibility = iota
	VisibilityFollowersOnly
	VisibilityNoOne
)

type SetLocationVisibilityInput struct {
	Visibility LocationVisibility `json:"visibility"`
}
