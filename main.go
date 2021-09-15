package main

import (
	"log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/nazeemnato/employee-go/routes"
)

func main() {
    app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))


    routes.Setup(app)
	
    log.Fatal(app.Listen(":8000"))
}