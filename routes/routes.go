package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazeemnato/employee-go/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("api/v1")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
}
