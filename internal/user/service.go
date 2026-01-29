package user

type UserService interface {
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
    Delete(id uint) error
    Update(id uint, input UpdateuserInput)
}