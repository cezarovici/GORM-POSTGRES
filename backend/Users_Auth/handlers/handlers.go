package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello Cezar from handlers")
}

// func CreateUser(c *fiber.Ctx) error {
// 	user := new(models.Users)

// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	database.DB.Db.Create(&user)

// 	return c.Status(200).JSON(user)
// }

// func Users(c *fiber.Ctx) error {
// 	users := []models.Users{}

// 	database.DB.Db.Find(&users)

// 	return c.Status(200).JSON(users)
// }
