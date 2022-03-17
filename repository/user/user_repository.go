package userrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/shopee/models"
	"github.com/praveennagaraj97/shopee/models/dto"
	"github.com/praveennagaraj97/shopee/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

// Method to initialize user repository
func (r *UserRepository) InitializeRepository(colln *mongo.Collection) {
	r.collection = colln

	// Index Options
	utils.CreateIndex(colln, bson.D{{Key: "email", Value: 1}}, "Email", true)
}

// Create / Save new user to database
func (r *UserRepository) Create(payload *dto.CreateUserDTO, role string) (*models.User, error) {
	var err error

	// create ctx for opertaions.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	err = payload.SetData(role)

	if err != nil {
		return nil, err
	}

	exist, err := r.checkUserExistsWithEmail(payload.Email)

	if err != nil {
		return nil, err
	}
	if exist {
		return nil, createError("email_already_taken")
	}
	_, err = r.collection.InsertOne(ctx, payload)

	if err != nil {
		return nil, createError("something_went_wrong")
	}

	return &models.User{
		ID:        payload.ID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
		JoinedOn:  payload.JoinedOn,
		UserRole:  payload.UserRole,
	}, nil
}

// Checks Email uniqueness against db.
func (r *UserRepository) checkUserExistsWithEmail(email string) (bool, error) {
	// create ctx for opertaions.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// check if user exists with email
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createError(msg string) error {
	return errors.New(msg)
}

// Finds the users by provided email and returs user entity.
func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result := r.collection.FindOne(ctx, bson.M{"email": email})

	if result.Err() != nil {
		return nil, createError("no_user_found_with_that_email")
	}

	var user models.User

	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Finds the users by provided email and returs user entity.
func (r *UserRepository) FindUserByID(id *primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result := r.collection.FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		return nil, errors.New("no_results_found")
	}

	var user models.User

	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update refresh token to user entity and return error if something goes wrong.
func (r *UserRepository) UpdateByField(uId primitive.ObjectID, key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := r.collection.UpdateByID(ctx, uId, bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}})
	if err != nil {
		return err
	}
	return nil
}

// Activate user marks user as active and saves refresh token.
func (r *UserRepository) ActivateUserAccount(uId primitive.ObjectID, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := r.collection.UpdateByID(ctx, uId, bson.D{{Key: "$set",
		Value: bson.D{{Key: "refresh_token", Value: refreshToken},
			{Key: "email_verified", Value: true},
			{Key: "is_active", Value: true}}}})
	if err != nil {
		return err
	}
	return nil
}

// Update User Data. returns Modified count
func (r *UserRepository) UpdateUser(id *primitive.ObjectID, data *dto.UpdateUserDTO) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	result, err := r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: data}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
