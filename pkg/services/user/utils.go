package user

import (
	"context"
	"sneaker-shop/pkg/database/models"
	"errors"
)

func ValidateAccount(
	login func(ctx context.Context, userID uint, userName string) (string, error),
	accountExists func(ctx context.Context, email string) (bool, uint, string, error),
	validatePassword func(ctx context.Context, user *models.LoginUser) (bool, error),
) func(ctx context.Context, user *models.LoginUser) (string, error) {
	return func(ctx context.Context, user *models.LoginUser) (string, error) {
		exists, userId, name, err := accountExists(ctx, user.Email)
		if err != nil {
			return "", err
		}
		if !exists {
			return "", errors.New("user not found")
		}

		_, err = validatePassword(ctx, user)
		if err != nil {
			return "", err
		}

		return login(ctx, userId, name)
	}
}
