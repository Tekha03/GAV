package vaccination

import (
	"time"

	"github.com/google/uuid"
)


type Vaccination struct {
	ID        uuid.UUID
	DogID     uuid.UUID
	Name      string
	DoneAt    time.Time
	NextDueAt *time.Time
	Notes     string
}
