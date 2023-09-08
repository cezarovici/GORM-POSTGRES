package main

import (
	"os"

	"github.com/cezarovici/GORM-POSTGRES/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	err := database.ConnectDb()
	if err != nil {
		os.Exit(1)
	}

	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
