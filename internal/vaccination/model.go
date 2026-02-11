package vaccination

import "time"


type Vaccination struct {
	ID        uint 
	DogID     uint   
	Name      string
	DoneAt    time.Time
	NextDueAt *time.Time
	Notes     string
}