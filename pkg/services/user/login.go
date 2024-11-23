package user

import (
	"context"
	"sneaker-shop/pkg/middleware"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
	"sneaker-shop/pkg/database/models"
	"fmt"
)

func (s *Service) Login(_ context.Context, userID uint, name string) (string, error) {
	now := time.Now()
	claims := middleware.UserClaims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "sneaker-shop",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *Service) Exists(_ context.Context, email string) (bool, uint, string, error) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return false, 0, "", err
	}

	return true, user.ID, user.Name, nil
}

func (s *Service) ValidatePassword(ctx context.Context, userInput *models.LoginUser) (bool, error) {
	user, err := s.getUserByEmail(userInput.Email)
	if err != nil {
		return false, err
	}

	err = userInput.CheckPassword(user.Password)
	if err != nil {
		return false, fmt.Errorf("invalid password: %w", err)
	}

	return true, nil
}

func (s *Service) getUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.db.Where(&models.User{Email: email}).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
