package helpers

import (
	"os"
	"time"

	"github.com/cezarovici/GORM-POSTGRES/models"

	db "github.com/cezarovici/GORM-POSTGRES/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(uuid string) (string, string) {
	claim, accessToken := GenerateAccessClaims(uuid)
	refreshToken := GenerateRefreshClaims(claim)

	return accessToken, refreshToken
}

// GenerateAccessClaims returns a claim and a acess_token string
func GenerateAccessClaims(uuid string) (*models.Claims, string) {
	t := time.Now()

	claim := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    uuid,
			ExpiresAt: &jwt.NumericDate{t.Add(time.Hour * 24)},
			Subject:   "access_token",
			IssuedAt:  &jwt.NumericDate{t.Add(time.Second * 1)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(cl *models.Claims) string {
	result := db.DB.Where(&models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: cl.Issuer,
		},
	}).Find(&models.Claims{})

	// checking the number of refresh tokens stored.
	// If the number is higher than 3, remove all the refresh tokens and leave only new one.
	if result.RowsAffected > 3 {
		db.DB.Where(&models.Claims{
			RegisteredClaims: jwt.RegisteredClaims{Issuer: cl.Issuer},
		}).Delete(&models.Claims{})
	}

	t := time.Now()
	refreshClaim := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: &jwt.NumericDate{t.Add(time.Hour * 24 * 30)},
			Subject:   "access_token",
			//IssuedAt:  t.Unix(),
		},
	}

	// create a claim on DB
	db.DB.Create(&refreshClaim)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}

// SecureAuth returns a middleware which secures all the private routes
// SecureAuth returns a middleware which secures all the private routes
func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		claims := new(models.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// Token signature is invalid
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			// Other error occurred
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !token.Valid {
			// Token is not valid
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			// Token has expired
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"general": "Token Expired",
			})
		}

		c.Locals("id", claims.Issuer)
		return c.Next()
	}
}

func GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	return accessCookie, refreshCookie
}
