package user

import (
	"context"
	"math"
	"time"

	"social_network/internal/dog"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

const (
	minNearbyRadiusMeters = 50
	maxNearbyRadiusMeters = 10_000
	walkingLocationTTL    = 2 * time.Hour
)

func NewService(repo Repository) (UserService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
}

func (s *service) Create(ctx context.Context, email, passwordHash string) (*User, error) {
	user, err := NewUser(email, passwordHash)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) error {
	user := &User{
		ID:       id,
		Email:    *input.Email,
		Password: *input.Password,
		Role:     *input.Role,
	}

	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindDogsNearby(ctx context.Context, userID uuid.UUID, centerLat, centerLon float64, radiusMeters float64) ([]*dog.Dog, error) {
	if err := validateCoordinates(centerLat, centerLon); err != nil {
		return nil, err
	}
	if err := validateRadius(radiusMeters); err != nil {
		return nil, err
	}

	dogs, err := s.repo.FindWalkingNearby(ctx, centerLat, centerLon, radiusMeters, time.Now().Add(-walkingLocationTTL))
	if err != nil {
		return nil, err
	}

	result := make([]*dog.Dog, 0)

	for _, currentDog := range dogs {
		owner, err := s.GetByID(ctx, currentDog.OwnerID)
		if err != nil {
			return nil, err
		}

		if currentDog.OwnerID == userID || owner.Visibility != VisibilityEveryone {
			continue
		}

		result = append(result, currentDog)
	}

	return result, nil
}

func (s *service) UpdateLocation(ctx context.Context, userID uuid.UUID, locationInput UpdateLocationInput) error {
	return s.withinTransaction(ctx, func(ctx context.Context) error {
		return s.updateLocation(ctx, userID, locationInput)
	})
}

func (s *service) updateLocation(ctx context.Context, userID uuid.UUID, locationInput UpdateLocationInput) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !locationInput.ClearLocation {
		if err := validateCoordinates(locationInput.Latitude, locationInput.Longitude); err != nil {
			return err
		}
		if err := validateLocationStatus(locationInput.Status); err != nil {
			return err
		}
		if err := validateLocationVisibility(locationInput.Visibility); err != nil {
			return err
		}
	}

	user.LocationStatus = locationInput.Status
	user.Visibility = locationInput.Visibility
	now := time.Now()
	user.LocationUpdatedAt = nil

	if locationInput.ClearLocation || locationInput.Status != Walking {
		user.Lat = nil
		user.Lon = nil
		user.LocationStatus = Inactive
		if locationInput.ClearLocation {
			user.Visibility = VisibilityNoOne
		}
		if err := s.repo.EndActiveWalkSession(ctx, userID, now); err != nil {
			return err
		}
		return s.repo.Update(ctx, user)
	}

	if err := s.repo.UpsertActiveWalkSession(ctx, &WalkSession{
		ID:         uuid.New(),
		UserID:     userID,
		Lat:        locationInput.Latitude,
		Lon:        locationInput.Longitude,
		Visibility: locationInput.Visibility,
		StartedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *service) SetLocationVisibility(ctx context.Context, userID uuid.UUID, visibility SetLocationVisibilityInput) error {
	return s.withinTransaction(ctx, func(ctx context.Context) error {
		return s.setLocationVisibility(ctx, userID, visibility)
	})
}

func (s *service) setLocationVisibility(ctx context.Context, userID uuid.UUID, visibility SetLocationVisibilityInput) error {
	if err := validateLocationVisibility(visibility.Visibility); err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Visibility = visibility.Visibility
	if err := s.repo.UpdateActiveWalkVisibility(ctx, userID, visibility.Visibility); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *service) withinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if transactional, ok := s.repo.(interface {
		WithinTransaction(context.Context, func(context.Context) error) error
	}); ok {
		return transactional.WithinTransaction(ctx, fn)
	}
	return fn(ctx)
}

func validateCoordinates(lat, lon float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 {
		return ErrInvalidLatitude
	}
	if math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 {
		return ErrInvalidLongitude
	}
	return nil
}

func validateRadius(radiusMeters float64) error {
	if math.IsNaN(radiusMeters) || math.IsInf(radiusMeters, 0) ||
		radiusMeters < minNearbyRadiusMeters || radiusMeters > maxNearbyRadiusMeters {
		return ErrInvalidRadius
	}
	return nil
}

func validateLocationStatus(status LocationStatus) error {
	switch status {
	case Inactive, Walking, ForcedOffline:
		return nil
	default:
		return ErrInvalidLocationStatus
	}
}

func validateLocationVisibility(visibility LocationVisibility) error {
	switch visibility {
	case VisibilityEveryone, VisibilityFollowersOnly, VisibilityNoOne:
		return nil
	default:
		return ErrInvalidLocationVisibility
	}
}
