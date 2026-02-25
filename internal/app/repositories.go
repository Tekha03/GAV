package app

import (
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
	gavSqlite "gav/storage/sqlite"

	"gorm.io/gorm"
)

type Repositories struct {
	User		user.Repository
	Profile		 profile.Repository
	Post		post.Repository
	Comment 	comment.Repository
	Like 		like.Repository
	Follow 		follow.Repository
	Dog 		dog.Repository
	Vaccination vaccination.Repository
	Stats 		stats.Repository
	Settings 	settings.Repository
}

func initRepositories(db *gorm.DB) (*Repositories, error) {
	r := &Repositories{}

	var err error
	r.User, 		err = gavSqlite.NewUserRepository(db);			if err != nil { return nil, err }
	r.Profile, 	 	 err = gavSqlite.NewProfileRepository(db);	 	  if err != nil { return nil, err }
	r.Post, 		err = gavSqlite.NewPostRepository(db);			if err != nil { return nil, err }
	r.Comment, 		err = gavSqlite.NewCommentRepository(db);		if err != nil { return nil, err }
	r.Like, 		err = gavSqlite.NewLikeRepository(db);			if err != nil { return nil, err }
	r.Follow, 		err = gavSqlite.NewFollowRepository(db);		if err != nil { return nil, err }
	r.Dog, 			err = gavSqlite.NewDogRepository(db);			if err != nil { return nil, err }
	r.Vaccination, 	err = gavSqlite.NewVaccinationRepository(db);	if err != nil { return nil, err }
	r.Stats, 		err = gavSqlite.NewStatsRepository(db);			if err != nil { return nil, err }
	r.Settings, 	err = gavSqlite.NewSettingsRepository(db);		if err != nil { return nil, err }

	return r, nil
}
