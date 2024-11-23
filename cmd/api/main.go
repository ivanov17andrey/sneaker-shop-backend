package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"flag"
	"github.com/rs/zerolog"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"sneaker-shop/pkg/database/models"
	"github.com/go-playground/validator/v10"
	"sneaker-shop/pkg/services/user"
	"sneaker-shop/pkg/handlers"
	user_handler "sneaker-shop/pkg/handlers/user"
)

var DB *gorm.DB

func main() {
	logLevel := flag.Int("log-level", 0, "Set log level (0 - debug, 1 - info, 2 - warn, 3 - error, 4 - fatal, 5 - panic)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	zerolog.SetGlobalLevel(zerolog.Level(*logLevel))
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Int("log-level", *logLevel).Msg("Logger initialized")

	// connect NATS
	// natsConn, err := nats.Connect(os.Getenv("NATS_SERVER_URL"), nats.Name("sneaker-shop"))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Error connecting to NATS")
	// }

	// if err := natsConn.Publish("orders", []byte("hello")); err != nil {
	// 	log.Fatal().Err(err).Msg("Error publishing message")
	// }

	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}

	server := handlers.NewServer(DB)

	validate := validator.New()

	// middlewares := []gin.HandlerFunc{middleware.AuthMiddleware()}

	userService := user.NewService(DB)
	user_handler.NewUserHandler(server, "/user", userService, validate)

	server.Gin.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	for _, route := range server.Gin.Routes() {
		log.Info().Str("method", route.Method).Str("path", route.Path).Msg("Route")
	}

	log.Fatal().Err(server.Run()).Msg("Error running server")
}
