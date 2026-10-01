package walk

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"social_network/internal/user"

	"github.com/google/uuid"
)

type repositoryStub struct {
	results    []Result
	err        error
	called     string
	userID     uuid.UUID
	latitude   float64
	longitude  float64
	radius     float64
	visibility user.LocationVisibility
	at         time.Time
}

func (r *repositoryStub) Start(_ context.Context, id uuid.UUID, lat, lon float64, visibility user.LocationVisibility, at time.Time) (*Result, error) {
	r.called, r.userID, r.latitude, r.longitude, r.visibility, r.at = "start", id, lat, lon, visibility, at
	return &Result{UserID: id, Latitude: lat, Longitude: lon, Visibility: visibility}, r.err
}
func (r *repositoryStub) Current(_ context.Context, id uuid.UUID, at time.Time) (*Result, error) {
	r.called, r.userID, r.at = "current", id, at
	return &Result{UserID: id}, r.err
}
func (r *repositoryStub) UpdateLocation(_ context.Context, id uuid.UUID, lat, lon float64, at time.Time) (*Result, error) {
	r.called, r.userID, r.latitude, r.longitude, r.at = "location", id, lat, lon, at
	return &Result{UserID: id, Latitude: lat, Longitude: lon}, r.err
}
func (r *repositoryStub) UpdateVisibility(_ context.Context, id uuid.UUID, visibility user.LocationVisibility, at time.Time) (*Result, error) {
	r.called, r.userID, r.visibility, r.at = "visibility", id, visibility, at
	return &Result{UserID: id, Visibility: visibility}, r.err
}
func (r *repositoryStub) Stop(_ context.Context, id uuid.UUID, at time.Time) error {
	r.called, r.userID, r.at = "stop", id, at
	return r.err
}
func (r *repositoryStub) Nearby(_ context.Context, id uuid.UUID, lat, lon, radius float64, at time.Time) ([]Result, error) {
	r.called, r.userID, r.latitude, r.longitude, r.radius, r.at = "nearby", id, lat, lon, radius, at
	return r.results, r.err
}

func fixedService(repo Repository) (*Service, time.Time) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	service := NewService(repo)
	service.now = func() time.Time { return now }
	return service, now
}

func TestServiceWalkLifecycleUsesExpectedTimes(t *testing.T) {
	repo := &repositoryStub{}
	service, now := fixedService(repo)
	ctx, userID := context.Background(), uuid.New()

	if _, err := service.Start(ctx, userID, 55.75, 37.62, user.VisibilityEveryone); err != nil {
		t.Fatal(err)
	}
	if repo.called != "start" || repo.at != now {
		t.Fatalf("Start() call = %q at %v", repo.called, repo.at)
	}
	if _, err := service.Current(ctx, userID); err != nil {
		t.Fatal(err)
	}
	if repo.called != "current" || repo.at != now.Add(-TTL) {
		t.Fatalf("Current() cutoff = %v, want %v", repo.at, now.Add(-TTL))
	}
	if _, err := service.UpdateLocation(ctx, userID, 55.76, 37.63); err != nil {
		t.Fatal(err)
	}
	if repo.called != "location" || repo.latitude != 55.76 || repo.longitude != 37.63 {
		t.Fatalf("UpdateLocation() captured %+v", repo)
	}
	if _, err := service.UpdateVisibility(ctx, userID, user.VisibilityFollowersOnly); err != nil {
		t.Fatal(err)
	}
	if repo.called != "visibility" || repo.visibility != user.VisibilityFollowersOnly {
		t.Fatalf("UpdateVisibility() captured %+v", repo)
	}
	if err := service.Stop(ctx, userID); err != nil {
		t.Fatal(err)
	}
	if repo.called != "stop" || repo.at != now {
		t.Fatalf("Stop() call = %q at %v", repo.called, repo.at)
	}
}

func TestServiceValidatesCoordinatesVisibilityAndRadius(t *testing.T) {
	service, _ := fixedService(&repositoryStub{})
	ctx, userID := context.Background(), uuid.New()

	for _, coordinates := range [][2]float64{{91, 0}, {-91, 0}, {0, 181}, {0, -181}, {math.NaN(), 0}, {0, math.Inf(1)}} {
		if _, err := service.Start(ctx, userID, coordinates[0], coordinates[1], user.VisibilityEveryone); !errors.Is(err, ErrInvalidCoordinates) {
			t.Fatalf("Start(%v) error = %v", coordinates, err)
		}
		if _, err := service.UpdateLocation(ctx, userID, coordinates[0], coordinates[1]); !errors.Is(err, ErrInvalidCoordinates) {
			t.Fatalf("UpdateLocation(%v) error = %v", coordinates, err)
		}
	}
	invalidVisibility := user.LocationVisibility(255)
	if _, err := service.Start(ctx, userID, 0, 0, invalidVisibility); !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("Start() error = %v", err)
	}
	if _, err := service.UpdateVisibility(ctx, userID, invalidVisibility); !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("UpdateVisibility() error = %v", err)
	}
	for _, radius := range []float64{49, 10001, math.NaN(), math.Inf(1)} {
		if _, err := service.Nearby(ctx, userID, 0, 0, radius); !errors.Is(err, ErrInvalidRadius) {
			t.Fatalf("Nearby(radius=%v) error = %v", radius, err)
		}
	}
}

func TestServiceNearbyFiltersSelfAndDistanceThenSorts(t *testing.T) {
	userID := uuid.New()
	nearID, nearerID, farID := uuid.New(), uuid.New(), uuid.New()
	repo := &repositoryStub{results: []Result{
		{UserID: userID, Latitude: 55.75, Longitude: 37.62},
		{UserID: nearID, Latitude: 55.751, Longitude: 37.62},
		{UserID: farID, Latitude: 56.0, Longitude: 37.62},
		{UserID: nearerID, Latitude: 55.7502, Longitude: 37.62},
	}}
	service, now := fixedService(repo)

	results, err := service.Nearby(context.Background(), userID, 55.75, 37.62, 1_000)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].UserID != nearerID || results[1].UserID != nearID {
		t.Fatalf("Nearby() results = %+v", results)
	}
	if results[0].DistanceMeters <= 0 || results[0].DistanceMeters >= results[1].DistanceMeters {
		t.Fatalf("distances not populated and sorted: %+v", results)
	}
	if repo.radius != 1_000 || repo.at != now.Add(-TTL) {
		t.Fatalf("repository query radius=%v cutoff=%v", repo.radius, repo.at)
	}
}

func TestServicePropagatesRepositoryErrors(t *testing.T) {
	wantErr := errors.New("repository failed")
	service, _ := fixedService(&repositoryStub{err: wantErr})
	if _, err := service.Nearby(context.Background(), uuid.New(), 0, 0, 100); !errors.Is(err, wantErr) {
		t.Fatalf("Nearby() error = %v, want %v", err, wantErr)
	}
}
