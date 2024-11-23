package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog/log"
	"github.com/go-playground/validator/v10"
	"regexp"
	"errors"
)

var emailRegex *regexp.Regexp

func init() {
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) HashPassword() {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal().Err(err).Msg("Error hashing password")
	}
	u.Password = string(hashedPwd)
}

func (l *LoginUser) CheckPassword(hashedPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(l.Password))
}

func NameValidator(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	return ok && str != ""
}

func EmailValidator(fl validator.FieldLevel) bool {
	email, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return emailRegex.MatchString(email)
}

func UserValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": "Unknown error"}
	}

	errorsMap := make(map[string]string)
	for _, e := range validationErrors {
		field := e.Field()
		switch e.Tag() {
		case "name":
			errorsMap[field] = "Provide your full name"
		case "email":
			errorsMap[field] = "Provide valid email address"
		default:
			errorsMap[field] = "Invalid"
		}
	}
	return errorsMap
}
