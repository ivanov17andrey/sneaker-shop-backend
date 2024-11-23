package user

import (
	"sneaker-shop/pkg/handlers"
	"sneaker-shop/pkg/database/models"
	"github.com/gin-gonic/gin"
	"sneaker-shop/pkg/services/user"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Handler struct {
	Server   *handlers.Server
	group    string
	router   *gin.RouterGroup
	service  *user.Service
	validate *validator.Validate
}

func NewUserHandler(s *handlers.Server, group string, service *user.Service, validate *validator.Validate) *Handler {
	handler := &Handler{
		Server:   s,
		group:    group,
		router:   s.Gin.Group(group),
		service:  service,
		validate: validate,
	}

	handler.routes()
	handler.registerValidator()

	return handler
}

func (h *Handler) registerValidator() {
	_ = h.validate.RegisterValidation("name", models.NameValidator)
	_ = h.validate.RegisterValidation("email", models.EmailValidator)
}

func (h *Handler) routes() http.Handler {
	h.router.POST("/", h.create)
	h.router.DELETE("/:id", h.delete)
	h.router.POST("/login", h.login)
	return h.Server.Gin
}
