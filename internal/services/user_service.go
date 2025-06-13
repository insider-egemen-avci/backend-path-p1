package services

import (
	"context"
	"fmt"
	"time"

	"insider-egemen-avci/backend-path-p1/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository models.UserRepository
}

func NewUserService(userRepository models.UserRepository) models.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email and password are required")
	}

	existingUser, err := service.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = service.userRepository.CreateUser(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (service *userService) Login(ctx context.Context, email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := service.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
