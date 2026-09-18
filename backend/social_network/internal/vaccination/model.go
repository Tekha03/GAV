package vaccination

import (
	"time"

	"github.com/google/uuid"
)

type Vaccination struct {
	ID        uuid.UUID  `json:"id"`
	DogID     uuid.UUID  `json:"dog_id"`
	Name      string     `json:"name"`
	DoneAt    time.Time  `json:"done_at"`
	NextDueAt *time.Time `json:"next_due_at"`
	Notes     string     `json:"notes"`
}
