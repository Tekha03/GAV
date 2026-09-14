package sqlite

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) (user.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &UserRepository{BaseRepository: repo}, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	var existing user.User
	err := r.DB(ctx).Where("email = ?", u.Email).First(&existing).Error

	if err == nil {
		return ErrUserExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.DB(ctx).Create(u).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	if err := r.DB(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.DB(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	updated := r.DB(ctx).Save(u)
	if updated.Error != nil {
		return updated.Error
	}

	if updated.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	deleted := r.DB(ctx).Delete(&user.User{}, "id = ?", id)
	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) FindWalkingNearby(ctx context.Context, centerLat, centerLon float64, radiusMeters float64, activeAfter time.Time) ([]*dog.Dog, error) {
	var dogs []*dog.Dog
	latDelta := radiusMeters / 111_320
	lonMeters := 111_320 * math.Cos(centerLat*math.Pi/180)
	if math.Abs(lonMeters) < 1 {
		lonMeters = 1
	}
	lonDelta := radiusMeters / lonMeters

	query := `
	SELECT
		dogs.id,
		dogs.owner_id,
		dogs.name,
		dogs.breed,
		dogs.photo_url,
		dogs.status,
		dogs.age,
		dogs.gender,
		walk_sessions.lat AS lat,
		walk_sessions.lon AS lon
	FROM dogs
	JOIN walk_sessions ON walk_sessions.user_id = dogs.owner_id
	WHERE walk_sessions.ended_at IS NULL
	  AND walk_sessions.visibility = ?
	  AND walk_sessions.updated_at >= ?
	  AND walk_sessions.lat BETWEEN ? AND ?
	  AND walk_sessions.lon BETWEEN ? AND ?
	`

	if err := r.DB(ctx).Raw(
		query,
		user.VisibilityEveryone,
		activeAfter,
		centerLat-latDelta,
		centerLat+latDelta,
		centerLon-lonDelta,
		centerLon+lonDelta,
	).Scan(&dogs).Error; err != nil {
		return nil, fmt.Errorf("dog repository: find walking nearby: %w", err)
	}

	result := make([]*dog.Dog, 0, len(dogs))
	for _, d := range dogs {
		if d.Lat == nil || d.Lon == nil {
			continue
		}
		distance := haversineMeters(centerLat, centerLon, *d.Lat, *d.Lon)
		if distance > radiusMeters {
			continue
		}
		d.DistanceMeters = &distance
		result = append(result, d)
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].DistanceMeters == nil {
			return false
		}
		if result[j].DistanceMeters == nil {
			return true
		}
		return *result[i].DistanceMeters < *result[j].DistanceMeters
	})

	return result, nil
}

func (r *UserRepository) UpsertActiveWalkSession(ctx context.Context, session *user.WalkSession) error {
	var existing user.WalkSession
	err := r.DB(ctx).
		Where("user_id = ? AND ended_at IS NULL", session.UserID).
		First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.DB(ctx).Create(session).Error
		}
		return err
	}

	return r.DB(ctx).
		Model(&user.WalkSession{}).
		Where("id = ?", existing.ID).
		Updates(map[string]any{
			"lat":        session.Lat,
			"lon":        session.Lon,
			"visibility": session.Visibility,
			"updated_at": session.UpdatedAt,
		}).Error
}

func (r *UserRepository) EndActiveWalkSession(ctx context.Context, userID uuid.UUID, endedAt time.Time) error {
	return r.DB(ctx).
		Model(&user.WalkSession{}).
		Where("user_id = ? AND ended_at IS NULL", userID).
		Updates(map[string]any{
			"ended_at":   endedAt,
			"updated_at": endedAt,
		}).Error
}

func (r *UserRepository) UpdateActiveWalkVisibility(ctx context.Context, userID uuid.UUID, visibility user.LocationVisibility) error {
	return r.DB(ctx).
		Model(&user.WalkSession{}).
		Where("user_id = ? AND ended_at IS NULL", userID).
		Update("visibility", visibility).Error
}

func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6_371_000

	toRad := func(value float64) float64 {
		return value * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	lat1Rad := toRad(lat1)
	lat2Rad := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}
