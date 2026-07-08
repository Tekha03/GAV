package dog

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDog_Success(t *testing.T) {
	ownerID := uuid.New()

	dog, err := NewDog(
		ownerID,
		"Buddy",
		"Labrador",
		Male,
		StatusFriendly,
		AdultAge,
		"url",
		"",
	)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, dog.ID)
	assert.Equal(t, ownerID, dog.OwnerID)
	assert.Equal(t, "Buddy", dog.Name)
}

func TestNewDog_OwnerNil(t *testing.T) {
	_, err := NewDog(uuid.Nil, "a", "b", Male, StatusFriendly, AdultAge, "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrOwnerIDNil, err)
}

func TestNewDog_NameEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "", "b", Male, StatusFriendly, AdultAge, "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrNameEmpty, err)
}

func TestNewDog_BreedEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "a", "", Male, StatusFriendly, AdultAge, "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrBreedEmpty, err)
}

func TestNewDog_GenderEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "a", "b", "", StatusFriendly, AdultAge, "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrGenderEmpty, err)
}

func TestNewDog_StatusEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "a", "b", Male, "", AdultAge, "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrStatusEmpty, err)
}

func TestNewDog_AgeEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "a", "b", Male, StatusFriendly, "", "url", "")

	assert.Error(t, err)
	assert.Equal(t, ErrAgeEmpty, err)
}

func TestNewDog_PhotoEmpty(t *testing.T) {
	_, err := NewDog(uuid.New(), "a", "b", Male, StatusFriendly, AdultAge, "", "")

	assert.Error(t, err)
	assert.Equal(t, ErrPhotoURLEmpty, err)
}
