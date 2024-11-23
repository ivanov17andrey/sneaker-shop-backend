package user

import (
	"context"
	"sneaker-shop/pkg/database/models"
	"fmt"
)

func (s *Service) Create(_ context.Context, user *models.User) (bool, error) {
	user.HashPassword()
	err := s.db.Create(user).Error
	if err != nil {
		return false, fmt.Errorf("failed to create user: %w", err)
	}

	return true, nil
}
