package memory

import (
	"context"
	"testing"

	"social_network/internal/comment"
	"social_network/internal/dog"
	"social_network/internal/post"

	"github.com/google/uuid"
)

func TestDogRepositoryLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewDogRepository()
	ownerID := uuid.New()
	dogModel := &dog.Dog{OwnerID: ownerID, Name: "Луна"}

	if err := repo.Create(ctx, nil); err != ErrDogNil {
		t.Fatalf("Create(nil) error = %v", err)
	}
	if err := repo.Create(ctx, dogModel); err != nil || dogModel.ID == uuid.Nil {
		t.Fatalf("Create() dog=%+v error=%v", dogModel, err)
	}
	if err := repo.Create(ctx, dogModel); err != ErrDogExists {
		t.Fatalf("duplicate Create() error = %v", err)
	}
	loaded, err := repo.GetByID(ctx, dogModel.ID)
	if err != nil || loaded != dogModel {
		t.Fatalf("GetByID() = %+v, %v", loaded, err)
	}
	dogModel.Name = "Луна II"
	if err := repo.Update(ctx, dogModel); err != nil {
		t.Fatal(err)
	}
	if dogs, err := repo.GetByOwnerID(ctx, ownerID); err != nil || len(dogs) != 1 || dogs[0].Name != "Луна II" {
		t.Fatalf("GetByOwnerID() = %+v, %v", dogs, err)
	}
	if err := repo.Update(ctx, &dog.Dog{ID: uuid.New()}); err != ErrDogNotFound {
		t.Fatalf("Update(missing) error = %v", err)
	}
	if err := repo.Delete(ctx, dogModel.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, dogModel.ID); err != ErrDogNotFound {
		t.Fatalf("GetByID(deleted) error = %v", err)
	}
	if err := repo.Delete(ctx, dogModel.ID); err != ErrDogNotFound {
		t.Fatalf("Delete(missing) error = %v", err)
	}
}

func TestPostRepositoryLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewPostRepository()
	userID := uuid.New()
	postModel := &post.Post{ID: uuid.New(), UserID: userID, Content: "Первая прогулка"}

	if err := repo.Create(ctx, nil); err != ErrPostNil {
		t.Fatalf("Create(nil) error = %v", err)
	}
	if err := repo.Create(ctx, postModel); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, postModel); err != ErrPostExists {
		t.Fatalf("duplicate Create() error = %v", err)
	}
	loaded, err := repo.GetByID(ctx, postModel.ID)
	if err != nil || loaded != postModel {
		t.Fatalf("GetByID() = %+v, %v", loaded, err)
	}
	if posts, err := repo.ListByUser(ctx, userID); err != nil || len(posts) != 1 {
		t.Fatalf("ListByUser() = %+v, %v", posts, err)
	}
	if err := repo.Delete(ctx, postModel.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, postModel.ID); err != ErrPostNotFound {
		t.Fatalf("GetByID(deleted) error = %v", err)
	}
	if err := repo.Delete(ctx, postModel.ID); err != ErrPostNotFound {
		t.Fatalf("Delete(missing) error = %v", err)
	}
}

func TestCommentRepositoryLifecycleAndOwnership(t *testing.T) {
	ctx := context.Background()
	repo := NewCommentRepository()
	postID, userID := uuid.New(), uuid.New()
	commentModel := &comment.Comment{PostID: postID, UserID: userID, Content: "Красиво!"}

	if err := repo.Create(ctx, nil); err != ErrCommentNil {
		t.Fatalf("Create(nil) error = %v", err)
	}
	if err := repo.Create(ctx, commentModel); err != nil || commentModel.ID == uuid.Nil || commentModel.CreatedAt.IsZero() {
		t.Fatalf("Create() comment=%+v error=%v", commentModel, err)
	}
	if comments, err := repo.ListyPostID(ctx, postID); err != nil || len(comments) != 1 || comments[0].ID != commentModel.ID {
		t.Fatalf("ListyPostID() = %+v, %v", comments, err)
	}
	if err := repo.Delete(ctx, commentModel.ID, uuid.New()); err != ErrCommentNotFound {
		t.Fatalf("Delete(wrong owner) error = %v", err)
	}
	if err := repo.Delete(ctx, commentModel.ID, userID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, commentModel.ID, userID); err != ErrCommentNotFound {
		t.Fatalf("Delete(missing) error = %v", err)
	}
}
