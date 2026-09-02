package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oscar/field-task-api/internal/auth"
	"github.com/oscar/field-task-api/internal/response"
)

func main() {
	// Load variables from a .env file if present. It is optional: in real
	// deployments the values usually come from the OS environment instead.
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file found; reading configuration from the environment")
	}

	jwtSecret := []byte(os.Getenv("AUTH_JWT_SECRET"))
	if len(jwtSecret) < 32 {
		log.Fatal("AUTH_JWT_SECRET must be at least 32 characters")
	}

	seedPassword := os.Getenv("AUTH_SEED_PASSWORD")
	if seedPassword == "" {
		seedPassword = "password123"
		log.Print("AUTH_SEED_PASSWORD is not set; using the development seed password")
	}

	authService, err := auth.NewService(seedPassword)
	if err != nil {
		log.Fatalf("initialize auth service: %v", err)
	}
	authHandler := auth.NewHandler(authService, jwtSecret)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/auth/login", authHandler.Login)
	router.GET("/auth/me", auth.RequireAuth(jwtSecret), authHandler.Me)

	server := &http.Server{
		Addr:              ":" + port(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func port() string {
	if value := os.Getenv("PORT"); value != "" {
		return value
	}

	return "8080"
}
