package walk

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
)

var (
	ErrInvalidCoordinates = errors.New("invalid coordinates")
	ErrInvalidVisibility  = errors.New("invalid visibility")
	ErrInvalidRadius      = errors.New("radius must be between 50 and 10000 meters")
	ErrNotFound           = errors.New("active walk not found")
	ErrAlreadyActive      = errors.New("walk already active")
)

const TTL = 2 * time.Hour

type Dog struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Breed    string     `json:"breed"`
	PhotoURL string     `json:"photo_url"`
	Notes    string     `json:"notes"`
	Status   dog.Status `json:"status"`
	Age      dog.Age    `json:"age"`
	Gender   dog.Gender `json:"gender"`
}

type Result struct {
	ID             uuid.UUID               `json:"id"`
	UserID         uuid.UUID               `json:"user_id"`
	Latitude       float64                 `json:"latitude"`
	Longitude      float64                 `json:"longitude"`
	Visibility     user.LocationVisibility `json:"visibility"`
	StartedAt      time.Time               `json:"started_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	DistanceMeters float64                 `json:"distance_meters,omitempty"`
	Dogs           []Dog                   `json:"dogs"`
}

type Repository interface {
	Start(context.Context, uuid.UUID, float64, float64, user.LocationVisibility, time.Time) (*Result, error)
	Current(context.Context, uuid.UUID, time.Time) (*Result, error)
	UpdateLocation(context.Context, uuid.UUID, float64, float64, time.Time) (*Result, error)
	UpdateVisibility(context.Context, uuid.UUID, user.LocationVisibility, time.Time) (*Result, error)
	Stop(context.Context, uuid.UUID, time.Time) error
	Nearby(context.Context, uuid.UUID, float64, float64, float64, time.Time) ([]Result, error)
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service { return &Service{repo: repo, now: time.Now} }

func (s *Service) Start(ctx context.Context, id uuid.UUID, lat, lon float64, visibility user.LocationVisibility) (*Result, error) {
	if err := validateCoordinates(lat, lon); err != nil {
		return nil, err
	}
	if err := validateVisibility(visibility); err != nil {
		return nil, err
	}
	return s.repo.Start(ctx, id, lat, lon, visibility, s.now())
}

func (s *Service) Current(ctx context.Context, id uuid.UUID) (*Result, error) {
	return s.repo.Current(ctx, id, s.now().Add(-TTL))
}

func (s *Service) UpdateLocation(ctx context.Context, id uuid.UUID, lat, lon float64) (*Result, error) {
	if err := validateCoordinates(lat, lon); err != nil {
		return nil, err
	}
	return s.repo.UpdateLocation(ctx, id, lat, lon, s.now())
}

func (s *Service) UpdateVisibility(ctx context.Context, id uuid.UUID, visibility user.LocationVisibility) (*Result, error) {
	if err := validateVisibility(visibility); err != nil {
		return nil, err
	}
	return s.repo.UpdateVisibility(ctx, id, visibility, s.now())
}

func (s *Service) Stop(ctx context.Context, id uuid.UUID) error { return s.repo.Stop(ctx, id, s.now()) }

func (s *Service) Nearby(ctx context.Context, id uuid.UUID, lat, lon, radius float64) ([]Result, error) {
	if err := validateCoordinates(lat, lon); err != nil {
		return nil, err
	}
	if math.IsNaN(radius) || math.IsInf(radius, 0) || radius < 50 || radius > 10000 {
		return nil, ErrInvalidRadius
	}
	results, err := s.repo.Nearby(ctx, id, lat, lon, radius, s.now().Add(-TTL))
	if err != nil {
		return nil, err
	}
	filtered := make([]Result, 0, len(results))
	for _, result := range results {
		if result.UserID == id {
			continue
		}
		distance := haversine(lat, lon, result.Latitude, result.Longitude)
		if distance <= radius {
			result.DistanceMeters = distance
			filtered = append(filtered, result)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].DistanceMeters < filtered[j].DistanceMeters })
	return filtered, nil
}

func validateCoordinates(lat, lon float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 || math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 {
		return ErrInvalidCoordinates
	}
	return nil
}

func validateVisibility(v user.LocationVisibility) error {
	if v != user.VisibilityEveryone && v != user.VisibilityFollowersOnly && v != user.VisibilityNoOne {
		return ErrInvalidVisibility
	}
	return nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000
	radians := func(v float64) float64 { return v * math.Pi / 180 }
	dLat, dLon := radians(lat2-lat1), radians(lon2-lon1)
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(radians(lat1))*math.Cos(radians(lat2))*math.Pow(math.Sin(dLon/2), 2)
	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
