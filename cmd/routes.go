package main

import (
	"github.com/cezarovici/GORM-POSTGRES/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/users", handlers.Users)
	app.Post("/user", handlers.CreateUser)
}
