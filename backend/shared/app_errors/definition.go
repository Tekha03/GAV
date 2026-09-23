package apperrors

type Definition struct {
	Code     Code
	Category Category
}

func define(code Code, category Category) Definition {
	return Definition{
		Code:     code,
		Category: category,
	}
}
