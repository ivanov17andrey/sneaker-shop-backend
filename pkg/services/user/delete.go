package user

import (
	"context"
	"sneaker-shop/pkg/database/models"
	"fmt"
)

func (s *Service) Delete(_ context.Context, userID uint) (bool, error) {
	err := s.db.Delete(&models.User{}, userID).Error
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	return true, nil
}
