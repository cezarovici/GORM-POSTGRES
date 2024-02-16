package helpers

import (
	"regexp"
	"strings"

	"github.com/cezarovici/GORM-POSTGRES/models"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// init initializes the validator
func init() {
	validate = validator.New()
}

// IsEmpty checks if a string is empty
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// ValidateRegister func validates the body of user for registration
func ValidateRegister(u *models.User) *models.UserErrors {
	e := &models.UserErrors{}

	// Validate username
	if IsEmpty(u.Username) {
		e.Err, e.Username = true, "Must not be empty"
	}

	// Validate email
	if err := validate.Var(u.Email, "email"); err != nil {
		e.Err, e.Email = true, "Must be a valid email"
	}

	// Validate password
	re := regexp.MustCompile("\\d") // regex check for at least one integer in string
	if !(len(u.Password) >= 8 && re.MatchString(u.Password)) {
		e.Err, e.Password = true, "Length of password should be atleast 8 and it must contain at least one number"
	}

	return e
}
