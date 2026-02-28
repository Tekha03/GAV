package dbserver

import (
	"fmt"
	"log/slog"
	"time"

	"gav/internal/comment"
	"gav/internal/dog"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/post"
	"gav/internal/profile"
	"gav/internal/settings"
	"gav/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, logger *slog.Logger) error {

	var count int64
	db.Model(&user.User{}).Count(&count)
	if count > 0 {
		logger.Info("seed skipped: database already has users")
		return nil
	}

	logger.Info("starting database seeding")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash seed password: %w", err)
	}
	pass := string(hashedPassword)

	users := []user.User{
		{ID: uuid.New(), Email: "admin@doglover.app", Password: pass, Role: "admin", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "anna.labrador@gmail.com", Password: pass, Role: "user", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "max.husky@yandex.ru", Password: pass, Role: "user", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "sophie.corgi@mail.ru", Password: pass, Role: "user", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "dima.shepherd@bk.ru", Password: pass, Role: "user", CreatedAt: time.Now()},
	}

	userMap := make(map[string]uuid.UUID)
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
		userMap[user.Email] = user.ID
	}

	adminID := userMap["admin@doglover.app"]
	annaID := userMap["anna.labrador@gmail.com"]
	maxID := userMap["max.husky@yandex.ru"]
	sophieID := userMap["sophie.corgi@mail.ru"]
	dimaID := userMap["dima.shepherd@bk.ru"]

	logger.Info("created 5 users")

	profiles := []profile.UserProfile{
		{UserID: adminID, Name: "Админ", Surname: "Главный", Username: "admin_dog", Bio: "Люблю всех собак одинаково"},
		{UserID: annaID, Name: "Анна", Surname: "Смирнова", Username: "anna_lab", Bio: "Моя лабрадор — смысл жизни"},
		{UserID: maxID, Name: "Максим", Surname: "Иванов", Username: "max_husky", Bio: "Хаски — это стиль жизни"},
		{UserID: sophieID, Name: "София", Surname: "Петрова", Username: "sophie_corgi", Bio: "Кorgi power!"},
		{UserID: dimaID, Name: "Дмитрий", Surname: "Кузнецов", Username: "dima_gsd", Bio: "Немецкая овчарка — лучший друг"},
	}

	for _, profile := range profiles {
		if err := db.Create(&profile).Error; err != nil {
			return fmt.Errorf("failed to create profile for %s: %w", profile.Username, err)
		}
	}

	settingsList := []settings.UserSettings{
		{UserID: adminID, ProfilePrivacy: false, ShowLocation: true, AllowMessages: true},
		{UserID: annaID, ProfilePrivacy: true, ShowLocation: false, AllowMessages: true},
		{UserID: maxID, ProfilePrivacy: false, ShowLocation: true, AllowMessages: true},
		{UserID: sophieID, ProfilePrivacy: true, ShowLocation: false, AllowMessages: false},
		{UserID: dimaID, ProfilePrivacy: false, ShowLocation: true, AllowMessages: true},
	}

	for _, setting := range settingsList {
		db.Create(&setting)
	}

	dogs := []dog.Dog{
		// Анна
		{ID: uuid.New(), OwnerID: annaID, Name: "Бимка", Breed: "Лабрадор-ретривер", Gender: dog.Female, Status: dog.StatusFriendly, Age: dog.AdultAge, PhotoUrl: "https://example.com/lab1.jpg"},
		{ID: uuid.New(), OwnerID: annaID, Name: "Рокси", Breed: "Лабрадор", Gender: dog.Female, Status: dog.StatusCautious, Age: dog.PuppyAge, PhotoUrl: "https://example.com/lab-puppy.jpg"},

		// Максим
		{ID: uuid.New(), OwnerID: maxID, Name: "Айс", Breed: "Сибирский хаски", Gender: dog.Male, Status: dog.StatusFriendly, Age: dog.PuppyAge, PhotoUrl: "https://example.com/husky1.jpg"},

		// София
		{ID: uuid.New(), OwnerID: sophieID, Name: "Булка", Breed: "Вельш-корги пемброк", Gender: dog.Male, Status: dog.StatusFriendly, Age: dog.AdultAge, PhotoUrl: "https://example.com/corgi1.jpg"},

		// Дмитрий
		{ID: uuid.New(), OwnerID: dimaID, Name: "Рекс", Breed: "Немецкая овчарка", Gender: dog.Male, Status: dog.StatusAggressive, Age: dog.ElderleAge, PhotoUrl: "https://example.com/gsd1.jpg"},
	}

	dogMap := make(map[string]uuid.UUID)
	for _, dog := range dogs {
		if err := db.Create(&dog).Error; err != nil {
			return fmt.Errorf("failed to create dog %s: %w", dog.Name, err)
		}
		dogMap[dog.Name] = dog.ID
	}

	logger.Info("created 5 dogs")

	posts := []post.Post{
		{ID: uuid.New(), UserID: annaID, Content: "Бимка сегодня поймала свой хвост за 3 секунды! 🐶", CreatedAt: time.Now().Add(-24 * time.Hour)},
		{ID: uuid.New(), UserID: maxID, Content: "Айс требует прогулку в -25°C и это нормально 😅❄️", CreatedAt: time.Now().Add(-12 * time.Hour)},
		{ID: uuid.New(), UserID: sophieID, Content: "Булка изобрела новый способ спать на двух стульях одновременно", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: uuid.New(), UserID: dimaID, Content: "Рекс охраняет дом от голубей. Миссия выполнена на 100%", CreatedAt: time.Now()},
	}

	for i, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return fmt.Errorf("failed to create post #%d, %w", i + 1, err)
		}
	}

	follows := []follow.Follow{
		{FollowerID: annaID, FollowingID: maxID},
		{FollowerID: annaID, FollowingID: sophieID},
		{FollowerID: maxID, FollowingID: annaID},
		{FollowerID: sophieID, FollowingID: annaID},
		{FollowerID: dimaID, FollowingID: annaID},
		{FollowerID: dimaID, FollowingID: maxID},
	}

	for _, follow := range follows {
		if err := db.Create(&follow).Error; err != nil {
			return fmt.Errorf("failed to create follow: %w", err)
		}
	}

	likes := []like.Like{
		{UserID: maxID, PostID: posts[0].ID},
		{UserID: sophieID, PostID: posts[0].ID},
		{UserID: annaID, PostID: posts[1].ID},
	}

	for _, like := range likes {
		if err := db.Create(&like); err != nil {
			return fmt.Errorf("failed to create like: %w", err)
		}
	}

	comments := []comment.Comment{
		{ID: uuid.New(), PostID: posts[0].ID, UserID: maxID, Content: "Какая умница! Скоро будет чемпион по ловле хвоста 😄", CreatedAt: time.Now()},
		{ID: uuid.New(), PostID: posts[2].ID, UserID: annaID, Content: "Это уровень босса сна! 😂", CreatedAt: time.Now()},
	}

	for _, comment := range comments {
		if err := db.Create(&comment); err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	logger.Info("database seeding completed successfully")
	return nil
}
