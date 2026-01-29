package user

type UserService interface {
    Register(email, password string) (*User, error)
    Login(email, password string) (*User, error)
}