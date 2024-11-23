package user

import (
	"github.com/gin-gonic/gin"
	"context"
	"time"
	"sneaker-shop/pkg/database/models"
	"net/http"
	"strconv"
	"sneaker-shop/pkg/services/user"
	"github.com/rs/zerolog/log"
)

func (h *Handler) create(c *gin.Context) {
	log.Info().Msg("Creating user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.validate.Struct(user); err != nil {
		validationErr := models.UserValidationErrors(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validationErr})
	}

	_, err := h.service.Create(ctx, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *Handler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 64)

	_, err = h.service.Delete(ctx, uint(idUint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handler) login(c *gin.Context) {
	var loginUser models.LoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	login := user.ValidateAccount(h.service.Login, h.service.Exists, h.service.ValidatePassword)
	token, err := login(c.Request.Context(), &loginUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
