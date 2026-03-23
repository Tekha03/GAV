package dog

import "errors"

var (
	ErrOwnerIDNil		= errors.New("dog model: owner id is nil")
	ErrNameEmpty		= errors.New("dog model: name is empty")
	ErrBreedEmpty		= errors.New("dog model: breed is empty")
	ErrGenderEmpty		= errors.New("dog model: gender is empty")
	ErrStatusEmpty		= errors.New("dog model: status is empty")
	ErrAgeEmpty			= errors.New("dog model: age is empty")
	ErrPhotoURLEmpty	= errors.New("dog model: photo url is empty")
	ErrRepoNil			= errors.New("dog service: repo is nil")
)
