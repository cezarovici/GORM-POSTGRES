package routes

import (
	"log"
	"math/rand"

	"github.com/cezarovici/GORM-POSTGRES/database"
	"github.com/cezarovici/GORM-POSTGRES/models"

	"time"

	"os"

	"github.com/cezarovici/GORM-POSTGRES/helpers"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/signup", CreateUser)              // Sign Up a user
	USER.Post("/signin", LoginUser)               // Sign In a user
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(helpers.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/user", GetUserData)

}

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	// validate if the email, username and password are in correct format
	errors := helpers.ValidateRegister(u)
	if errors.Err {
		return c.JSON(errors)
	}

	log.Print(database.DB)

	if count := database.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Email is already registered"
	}
	if count := database.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Username = true, "Username is already registered"
	}
	if errors.Err {
		return c.JSON(errors)
	}

	// Hashing the password with a random salt
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		password,
		rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
	)

	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)

	if err := database.DB.Create(&u).Error; err != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := helpers.GenerateTokens(u.UUID.String())
	accessCookie, refreshCookie := helpers.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

// LoginUser route logins a user in the app
func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// check if a user exists
	u := new(models.User)
	if res := database.DB.Where(
		&models.User{Email: input.Identity}).Or(
		&models.User{Username: input.Identity},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := helpers.GenerateTokens(u.UUID.String())
	accessCookie, refreshCookie := helpers.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetUserData returns the details of the user signed in
func GetUserData(c *fiber.Ctx) error {
	id := c.Locals("id")

	u := new(models.User)
	if res := database.DB.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}

// GetAccessToken generates and sends a new access token iff there is a valid refresh token
func GetAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	refreshClaims := new(models.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := database.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&models.Claims{}); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt.Unix() < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken := helpers.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}
