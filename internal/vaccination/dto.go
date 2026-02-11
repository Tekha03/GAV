package vaccination

import "time"

type CreateVaccinationInput struct {
	Name      string
	DoneAt    time.Time
	NextDueAt *time.Time
	Notes     string
}

type UpdateVaccinationInput struct {
	Name      *string
	DoneAt    *time.Time
	NextDueAt *time.Time
	Notes     *string
}