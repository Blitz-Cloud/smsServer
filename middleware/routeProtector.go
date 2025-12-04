package middleware

import (
	"strings"

	"github.com/Blitz-Cloud/smsServer/utils"
	"github.com/gofiber/fiber/v2"
)

func RouteProtector(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if len(authorizationHeader) == 0 || len(authorizationHeader) <= len("Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid format for authorization header, or authorization header missing")
	}
	token := strings.Split(authorizationHeader, " ")[1]
	_, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	return c.Next()
}
