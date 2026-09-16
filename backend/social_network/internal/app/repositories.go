package app

import (
	"social_network/internal/comment"
	"social_network/internal/device"
	"social_network/internal/dog"
	"social_network/internal/follow"
	"social_network/internal/like"
	"social_network/internal/notification"
	"social_network/internal/post"
	"social_network/internal/profile"
	"social_network/internal/settings"
	"social_network/internal/stats"
	"social_network/internal/token"
	"social_network/internal/user"
	"social_network/internal/vaccination"
	"social_network/storage/postgres"

	"gorm.io/gorm"
)

type Repositories struct {
	User         user.Repository
	Token        token.Repository
	Profile      profile.Repository
	Post         post.Repository
	Comment      comment.Repository
	Like         like.Repository
	Follow       follow.Repository
	Dog          dog.Repository
	Vaccination  vaccination.Repository
	Stats        stats.Repository
	Settings     settings.Repository
	Notification notification.Repository
	Device       device.Repository
}

func initRepositories(db *gorm.DB) (*Repositories, error) {
	r := &Repositories{}

	var err error
	r.User, err = postgres.NewUserRepository(db)
	if err != nil {
		return nil, err
	}
	r.Token, err = postgres.NewTokenRepository(db)
	if err != nil {
		return nil, err
	}
	r.Profile, err = postgres.NewProfileRepository(db)
	if err != nil {
		return nil, err
	}
	r.Post, err = postgres.NewPostRepository(db)
	if err != nil {
		return nil, err
	}
	r.Comment, err = postgres.NewCommentRepository(db)
	if err != nil {
		return nil, err
	}
	r.Like, err = postgres.NewLikeRepository(db)
	if err != nil {
		return nil, err
	}
	r.Follow, err = postgres.NewFollowRepository(db)
	if err != nil {
		return nil, err
	}
	r.Dog, err = postgres.NewDogRepository(db)
	if err != nil {
		return nil, err
	}
	r.Vaccination, err = postgres.NewVaccinationRepository(db)
	if err != nil {
		return nil, err
	}
	r.Stats, err = postgres.NewStatsRepository(db)
	if err != nil {
		return nil, err
	}
	r.Settings, err = postgres.NewSettingsRepository(db)
	if err != nil {
		return nil, err
	}
	r.Notification, err = postgres.NewNotificationRepository(db)
	if err != nil {
		return nil, err
	}
	r.Device, err = postgres.NewDeviceRepo(db)
	if err != nil {
		return nil, err
	}

	return r, nil
}
