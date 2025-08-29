package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"userService/internal/model"
)

type MongoUserRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
		timeout:    5 * time.Second,
	}
}

// CreateUser inserts a new user with CreatedAt and UpdatedAt timestamps.
func (r *MongoUserRepository) CreateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user already exists")
	}
	return err
}

// UpdateUser updates mutable fields and sets UpdatedAt timestamp.
func (r *MongoUserRepository) UpdateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	user.UpdatedAt = time.Now().UTC()

	filter := bson.M{
		"_id":        user.ID,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{
		"firstname":     user.Firstname,
		"lastname":      user.Lastname,
		"email":         user.Email,
		"about":         user.About,
		"date_of_birth": user.DateOfBirth,
		"avatar_url":    user.AvatarURL,
		"gender":        user.Gender,
		"location":      user.Location,
		"socials":       user.Socials,
		"updated_at":    user.UpdatedAt,
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found or deleted")
	}
	return nil
}

// DeleteUserById performs a soft delete by setting DeletedAt timestamp.
func (r *MongoUserRepository) DeleteUserById(userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{
		"_id":        userId,
		"deleted_at": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{"deleted_at": time.Now().UTC()},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found or already deleted")
	}
	return nil
}

// GetAllUsers returns all non-deleted users.
func (r *MongoUserRepository) GetAllUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cur, err := r.collection.Find(ctx, bson.M{"deleted_at": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []model.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserById finds a non-deleted user by ID.
func (r *MongoUserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user model.User
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}
