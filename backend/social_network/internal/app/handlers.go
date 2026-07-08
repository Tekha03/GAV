package app

import (
	"social_network/internal/notification"
	"social_network/transport/http/handlers"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	User        *handlers.UserHandler
	Profile     *handlers.ProfileHandler
	Post        *handlers.PostHandler
	Feed        *handlers.FeedHandler
	Comment     *handlers.CommentHandler
	Like        *handlers.LikeHandler
	Follow      *handlers.FollowHandler
	Dog         *handlers.DogHandler
	Vaccination *handlers.VaccinationHandler
	Stats       *handlers.StatsHandler
	Settings    *handlers.SettingsHandler
	Upload      *handlers.UploadHandler
	WSHandler   *handlers.NotificationHandler
}

func initHandlers(services *Services, notificationHub *notification.Hub) (*Handlers, error) {
	h := &Handlers{}

	var err error
	h.Auth, err = handlers.NewAuthHandler(services.Auth)
	if err != nil {
		return nil, err
	}
	h.User, err = handlers.NewUserHandler(services.User)
	if err != nil {
		return nil, err
	}
	h.Profile, err = handlers.NewProfileHandler(services.Profile)
	if err != nil {
		return nil, err
	}
	h.Post, err = handlers.NewPostHandler(services.Post,
		services.Media)
	if err != nil {
		return nil, err
	}
	h.Feed, err = handlers.NewFeedHandler(services.Feed)
	if err != nil {
		return nil, err
	}
	h.Comment, err = handlers.NewCommentHandler(services.Comment,
		services.Post, services.Notification)
	if err != nil {
		return nil, err
	}
	h.Like, err = handlers.NewLikeHandler(services.Like,
		services.Post,
		services.Notification)
	if err != nil {
		return nil, err
	}
	h.Follow, err = handlers.NewFollowHandler(services.Follow,
		services.Notification)
	if err != nil {
		return nil, err
	}
	h.Dog, err = handlers.NewDogHandler(services.Dog)
	if err != nil {
		return nil, err
	}
	h.Vaccination, err = handlers.NewVaccinationHandler(services.Vaccination)
	if err != nil {
		return nil, err
	}
	h.Stats, err = handlers.NewStatsHandler(services.Stats)
	if err != nil {
		return nil, err
	}
	h.Settings, err = handlers.NewSettingsHandler(services.Settings)
	if err != nil {
		return nil, err
	}
	h.Upload, err = handlers.NewUploadHandler(services.Media)
	if err != nil {
		return nil, err
	}
	h.WSHandler, err = handlers.NewNotificationHandler(notificationHub)
	if err != nil {
		return nil, err
	}

	return h, nil
}
