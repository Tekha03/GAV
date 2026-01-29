package profile

type CreateProfileInput struct {
	Name         string
    Surname      string
    Username     string
    ProfilePhoto string
    Bio          string
    Address      string
    BirthDate    string
}

type UpdateProfileInput struct {
	Name         *string
    Surname      *string
    Username     *string
    ProfilePhoto *string
    Bio          *string
    Address      *string
    BirthDate    *string
}