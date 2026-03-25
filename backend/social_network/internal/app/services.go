package app

import (
	"social_network/internal/auth"
	"social_network/internal/comment"
	"social_network/internal/dog"
	"social_network/internal/feed"
	"social_network/internal/follow"
	"social_network/internal/like"
	"social_network/internal/media"
	"social_network/internal/notification"
	"social_network/internal/post"
	"social_network/internal/profile"
	"social_network/internal/settings"
	"social_network/internal/stats"
	"social_network/internal/token"
	"social_network/internal/user"
	"social_network/internal/vaccination"
)

type Services struct {
	User		user.UserService
	Token		token.TokenService
	Auth		auth.AuthService
	Profile		 profile.ProfileService
	Post		post.PostService
	Feed		feed.FeedService
	Media		media.MediaService
	Comment 	comment.CommentService
	Like 		like.LikeService
	Follow 		follow.FollowService
	Notification notification.NotificationService
	Dog 		dog.DogService
	Vaccination vaccination.VaccinationService
	Stats 		stats.StatsService
	Settings 	settings.SettingsService
}

func initServices(repos *Repositories, jwtConfig auth.JWTConfig, storage media.Storage, notificationHub *notification.Hub) (*Services, error) {
	s := &Services{}

	var err error
	s.User, 		err = user.NewService(repos.User);		 			if err != nil { return nil, err }
	s.Token,		err = token.NewService(repos.Token);				if err != nil { return nil, err }
	s.Auth, 		err = auth.NewService(s.User, s.Token, jwtConfig, &auth.PasswordHasher{});		if err != nil { return nil, err }
	s.Profile, 	 	 err = profile.NewService(repos.Profile);	  		   if err != nil { return nil, err }
	s.Post, 		err = post.NewService(repos.Post);					if err != nil { return nil, err }
	s.Feed,			err = feed.NewService(repos.Post);					if err != nil { return nil, err }
	s.Media, 		err = media.NewService(storage);					if err != nil { return nil, err }
	s.Comment, 		err = comment.NewService(repos.Comment); 			if err != nil { return nil, err }
	s.Like, 		err = like.NewService(repos.Like);					if err != nil { return nil, err }
	s.Follow, 		err = follow.NewService(repos.Follow);				if err != nil { return nil, err }
	s.Notification,	 err = notification.NewService(notificationHub, repos.Notification);	   if err != nil { return nil, err }
	s.Dog, 			err = dog.NewService(repos.Dog);					if err != nil { return nil, err }
	s.Vaccination, 	err = vaccination.NewService(repos.Vaccination);	if err != nil { return nil, err }
	s.Stats, 		err = stats.NewService(repos.Stats);				if err != nil { return nil, err }
	s.Settings, 	err = settings.NewService(repos.Settings);			if err != nil { return nil, err }

	return s, nil
}
