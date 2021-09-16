package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazeemnato/employee-go/controllers"
	"github.com/nazeemnato/employee-go/middlewares"
)

func Setup(app *fiber.App) {
	api := app.Group("api/v1")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	proctedApi := api.Use(middlewares.IsAuthenticated)
	proctedApi.Get("/user", controllers.User)
	proctedApi.Post("/logout", controllers.Logout)
}
