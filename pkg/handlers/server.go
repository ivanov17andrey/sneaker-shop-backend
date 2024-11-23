package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/cors"
	"github.com/rs/zerolog"
	"fmt"
	"os"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Gin *gin.Engine
	db  *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	ge := gin.New()

	ge.Use(logger.SetLogger(logger.WithLogger(func(*gin.Context, zerolog.Logger) zerolog.Logger {
		return log.Logger
	})))
	ge.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},                                       // List of allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // List of allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow the necessary headers
		AllowCredentials: true,
	}

	ge.Use(cors.New(corsConfig))

	return &Server{
		Gin: ge,
		db:  db,
	}
}

func (s *Server) Run() error {
	return s.Gin.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
