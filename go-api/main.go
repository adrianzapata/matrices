// Command matrices-go-api starts the Go/Fiber API described in the
// Interseguro coding challenge: it receives a matrix, rotates it and
// computes its QR factorization, forwards the results to the Node.js
// statistics API, and returns the combined response.
package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"

	"github.com/interseguro/matrices-go-api/client"
	"github.com/interseguro/matrices-go-api/handlers"
	"github.com/interseguro/matrices-go-api/middleware"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// mintServiceToken creates a long-lived HS256 JWT identifying the Go API as
// the caller, signed with the secret shared between both services.
func mintServiceToken(secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": "go-api",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func main() {
	port := getEnv("PORT", "8080")
	nodeAPIURL := getEnv("NODE_API_URL", "http://localhost:4000")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-me")
	authUsername := getEnv("AUTH_USERNAME", "admin")
	authPassword := getEnv("AUTH_PASSWORD", "admin")

	// Both APIs trust the same JWT_SECRET, so the Go API mints its own
	// long-lived service-to-service token at startup and uses it to
	// authenticate every call it makes to the Node.js stats API.
	serviceToken, err := mintServiceToken(jwtSecret)
	if err != nil {
		log.Fatalf("failed to mint service token: %v", err)
	}

	nodeClient := client.NewNodeClient(nodeAPIURL, serviceToken)

	app := fiber.New(fiber.Config{
		AppName: "matrices-go-api",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("CORS_ORIGINS", "*"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	authCfg := handlers.AuthConfig{
		Username: authUsername,
		Password: authPassword,
		Secret:   jwtSecret,
		TTL:      time.Hour,
	}

	api := app.Group("/api")
	api.Post("/auth/login", handlers.Login(authCfg))

	protected := api.Group("/matrix", middleware.JWTAuth(jwtSecret))
	protected.Post("/process", handlers.ProcessMatrix(nodeClient))

	log.Printf("matrices-go-api listening on :%s (node api: %s)", port, nodeAPIURL)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
