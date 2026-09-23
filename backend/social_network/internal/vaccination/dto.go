package vaccination

import "time"

type CreateVaccinationInput struct {
	Name      string     `json:"name"`
	DoneAt    time.Time  `json:"done_at"`
	NextDueAt *time.Time `json:"next_due_at"`
	Notes     string     `json:"notes"`
}

type UpdateVaccinationInput struct {
	Name      *string    `json:"name"`
	DoneAt    *time.Time `json:"done_at"`
	NextDueAt *time.Time `json:"next_due_at"`
	Notes     *string    `json:"notes"`
}
