package postgres

import (
	"context"
	"errors"
	"math"
	"time"

	"social_network/internal/dog"
	"social_network/internal/user"
	"social_network/internal/walk"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalkRepository struct{ db *gorm.DB }

func NewWalkRepository(db *gorm.DB) *WalkRepository { return &WalkRepository{db: db} }

func (r *WalkRepository) Start(ctx context.Context, userID uuid.UUID, lat, lon float64, visibility user.LocationVisibility, now time.Time) (*walk.Result, error) {
	var session user.WalkSession
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Serializes concurrent starts for the same user on PostgreSQL.
		var account user.User
		lockQuery := "SELECT * FROM users WHERE id = ?"
		if tx.Dialector.Name() == "postgres" {
			lockQuery += " FOR UPDATE"
		}
		if err := tx.Raw(lockQuery, userID).Scan(&account).Error; err != nil {
			return err
		}
		if account.ID == uuid.Nil {
			return user.ErrUserNotFound
		}
		if err := tx.Where("user_id = ? AND ended_at IS NULL", userID).First(&session).Error; err == nil {
			if session.UpdatedAt.After(now.Add(-walk.TTL)) {
				return walk.ErrAlreadyActive
			}
			if err := tx.Model(&session).Updates(map[string]any{"ended_at": now, "updated_at": now}).Error; err != nil {
				return err
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		session = user.WalkSession{ID: uuid.New(), UserID: userID, Lat: lat, Lon: lon, Visibility: visibility, StartedAt: now, UpdatedAt: now}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}
		return tx.Model(&user.User{}).Where("id = ?", userID).Updates(map[string]any{
			"location_status": user.Walking, "visibility": visibility, "location_updated_at": now,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return toWalkResult(session), nil
}

func (r *WalkRepository) Current(ctx context.Context, userID uuid.UUID, activeAfter time.Time) (*walk.Result, error) {
	var session user.WalkSession
	err := r.db.WithContext(ctx).Where("user_id = ? AND ended_at IS NULL AND updated_at >= ?", userID, activeAfter).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, walk.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toWalkResult(session), nil
}

func (r *WalkRepository) UpdateLocation(ctx context.Context, userID uuid.UUID, lat, lon float64, now time.Time) (*walk.Result, error) {
	return r.update(ctx, userID, now, map[string]any{"lat": lat, "lon": lon, "updated_at": now}, map[string]any{"location_updated_at": now})
}

func (r *WalkRepository) UpdateVisibility(ctx context.Context, userID uuid.UUID, visibility user.LocationVisibility, now time.Time) (*walk.Result, error) {
	return r.update(ctx, userID, now, map[string]any{"visibility": visibility, "updated_at": now}, map[string]any{"visibility": visibility, "location_updated_at": now})
}

func (r *WalkRepository) update(ctx context.Context, userID uuid.UUID, now time.Time, fields, userFields map[string]any) (*walk.Result, error) {
	var session user.WalkSession
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user.WalkSession{}).Where("user_id = ? AND ended_at IS NULL AND updated_at >= ?", userID, now.Add(-walk.TTL)).Updates(fields)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return walk.ErrNotFound
		}
		updatedUser := tx.Model(&user.User{}).Where("id = ?", userID).Updates(userFields)
		if updatedUser.Error != nil {
			return updatedUser.Error
		}
		if updatedUser.RowsAffected == 0 {
			return user.ErrUserNotFound
		}
		return tx.Where("user_id = ? AND ended_at IS NULL", userID).First(&session).Error
	})
	if err != nil {
		return nil, err
	}
	return toWalkResult(session), nil
}

func (r *WalkRepository) Stop(ctx context.Context, userID uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user.WalkSession{}).Where("user_id = ? AND ended_at IS NULL", userID).Updates(map[string]any{"ended_at": now, "updated_at": now})
		if result.Error != nil || result.RowsAffected == 0 {
			return result.Error
		}
		return tx.Model(&user.User{}).Where("id = ?", userID).Updates(map[string]any{"location_status": user.Inactive, "location_updated_at": nil}).Error
	})
}

func (r *WalkRepository) Nearby(ctx context.Context, requesterID uuid.UUID, lat, lon, radius float64, activeAfter time.Time) ([]walk.Result, error) {
	// Conservative bounds avoid dropping points near the radius edge before Haversine filtering.
	latDelta := radius / 110000
	lonMeters := 110000 * math.Cos(lat*math.Pi/180)
	if math.Abs(lonMeters) < 1 {
		lonMeters = 1
	}
	lonDelta := radius / lonMeters
	var sessions []user.WalkSession
	query := `
		SELECT ws.* FROM walk_sessions ws
		WHERE ws.ended_at IS NULL AND ws.updated_at >= ? AND ws.user_id <> ?
		AND ws.lat BETWEEN ? AND ?
		AND (ws.visibility = ? OR (ws.visibility = ? AND EXISTS (
			SELECT 1 FROM follows f WHERE f.follower_id = ? AND f.following_id = CAST(ws.user_id AS TEXT)
		)))`
	args := []any{activeAfter, requesterID, lat - latDelta, lat + latDelta, user.VisibilityEveryone, user.VisibilityFollowersOnly, requesterID}
	if lon-lonDelta >= -180 && lon+lonDelta <= 180 {
		query += " AND ws.lon BETWEEN ? AND ?"
		args = append(args, lon-lonDelta, lon+lonDelta)
	}
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&sessions).Error
	if err != nil {
		return nil, err
	}
	results := make([]walk.Result, 0, len(sessions))
	ids := make([]uuid.UUID, 0, len(sessions))
	indexes := make(map[uuid.UUID]int, len(sessions))
	for _, session := range sessions {
		if index, exists := indexes[session.UserID]; exists {
			if session.UpdatedAt.After(results[index].UpdatedAt) {
				results[index] = *toWalkResult(session)
			}
			continue
		}
		indexes[session.UserID] = len(results)
		results = append(results, *toWalkResult(session))
		ids = append(ids, session.UserID)
	}
	if len(ids) == 0 {
		return results, nil
	}
	var dogs []walkDogRow
	err = r.db.WithContext(ctx).Raw("SELECT id, owner_id, name, breed, photo_url, notes, status, age, gender FROM dogs WHERE owner_id IN ?", ids).Scan(&dogs).Error
	if err != nil {
		return nil, err
	}
	for _, dog := range dogs {
		index := indexes[dog.OwnerID]
		results[index].Dogs = append(results[index].Dogs, walk.Dog{ID: dog.ID, Name: dog.Name, Breed: dog.Breed, PhotoURL: dog.PhotoURL, Notes: dog.Notes, Status: dog.Status, Age: dog.Age, Gender: dog.Gender})
	}
	return results, nil
}

type walkDogRow struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	Name     string
	Breed    string
	PhotoURL string
	Notes    string
	Status   dog.Status
	Age      dog.Age
	Gender   dog.Gender
}

func toWalkResult(s user.WalkSession) *walk.Result {
	return &walk.Result{ID: s.ID, UserID: s.UserID, Latitude: s.Lat, Longitude: s.Lon, Visibility: s.Visibility, StartedAt: s.StartedAt, UpdatedAt: s.UpdatedAt, Dogs: []walk.Dog{}}
}
