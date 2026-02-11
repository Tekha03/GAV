package dbserver

type RegisterDTO struct {
	Email	 string	`json:"email" binding:"required,email"`
	Password string	`json:"password" binding:"required,min=6"`
}

type LoginDTO struct {
	Email	 string	`json:"email" binding:"required,email"`
	Password string	`json:"password" binding:"required"`
}
