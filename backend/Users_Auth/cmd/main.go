package main

import (
	"os"

	"github.com/cezarovici/GORM-POSTGRES/database"
	"github.com/cezarovici/GORM-POSTGRES/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	port = "3000"
)

func CreateServer() *fiber.App {
	return fiber.New()
}

func main() {
	err := database.ConnectDb()
	if err != nil {
		os.Exit(1)
	}

	app := CreateServer()
	app.Use(cors.New())

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
