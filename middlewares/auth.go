package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx)  error{ 
	cookie := c.Cookies("token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("sce"), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Invalid token",
		})
	}
	payload := token.Claims.(*jwt.StandardClaims)

	c.Append("x-userId", payload.Subject)

	return c.Next()

}