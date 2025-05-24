package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"dnd-simulator/internal/database"
	"dnd-simulator/internal/models"
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	collection := s.db.GetCollection("users")

	// Check if user already exists
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return nil, errors.New("user with this username or email already exists")
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = "" // Don't return password
	return user, nil
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	collection := s.db.GetCollection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": username},
		},
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = "" // Don't return password
	return &user, nil
}

func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	collection := s.db.GetCollection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	user.Password = "" // Don't return password
	return &user, nil
}