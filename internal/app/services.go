package app

import (
	"gav/internal/auth"
	"gav/internal/comment"
	"gav/internal/dog"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/post"
	"gav/internal/profile"
	"gav/internal/settings"
	"gav/internal/stats"
	"gav/internal/user"
	"gav/internal/vaccination"
)

type Services struct {
	User		user.UserService
	Auth		auth.AuthService
	Profile		 profile.ProfileService
	Post		post.PostService
	Comment 	comment.CommentService
	Like 		like.LikeService
	Follow 		follow.FollowService
	Dog 		dog.DogService
	Vaccination vaccination.VaccinationService
	Stats 		stats.StatsService
	Settings 	settings.SettingsService
}

func initServices(repos *Repositories, jwtConfig auth.JWTConfig) (*Services, error) {
	s := &Services{}

	var err error
	s.User, 		err = user.NewService(repos.User);		 			if err != nil { return nil, err }
	s.Auth, 		err = auth.NewService(s.User, jwtConfig, &auth.PasswordHasher{});		if err != nil { return nil, err }
	s.Profile, 	 	 err = profile.NewService(repos.Profile);	  		   if err != nil { return nil, err }
	s.Post, 		err = post.NewService(repos.Post);					if err != nil { return nil, err }
	s.Comment, 		err = comment.NewService(repos.Comment); 			if err != nil { return nil, err }
	s.Like, 		err = like.NewService(repos.Like);					if err != nil { return nil, err }
	s.Follow, 		err = follow.NewService(repos.Follow);				if err != nil { return nil, err }
	s.Dog, 			err = dog.NewService(repos.Dog);					if err != nil { return nil, err }
	s.Vaccination, 	err = vaccination.NewService(repos.Vaccination);	if err != nil { return nil, err }
	s.Stats, 		err = stats.NewService(repos.Stats);				if err != nil { return nil, err }
	s.Settings, 	err = settings.NewService(repos.Settings);			if err != nil { return nil, err }

	return s, nil
}
