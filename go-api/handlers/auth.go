package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthConfig holds the credentials and secret used to issue demo JWTs.
// This is intentionally simple (single hardcoded user via env vars) since
// authentication itself is not the focus of the challenge — it only exists
// to demonstrate protecting the matrix endpoints with JWT, per the
// "funcionalidad opcional" requirement.
type AuthConfig struct {
	Username string
	Password string
	Secret   string
	TTL      time.Duration
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresInSeconds"`
}

// Login issues a short-lived JWT when the provided credentials match cfg.
func Login(cfg AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		}

		if req.Username != cfg.Username || req.Password != cfg.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		claims := jwt.MapClaims{
			"sub": req.Username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(cfg.TTL).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(cfg.Secret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to sign token"})
		}

		return c.JSON(loginResponse{Token: signed, ExpiresIn: int(cfg.TTL.Seconds())})
	}
}
