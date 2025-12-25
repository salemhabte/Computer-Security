package infrastructure

import (
	"errors"
	"fmt"
	"regexp"
	"security/config"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}
func (p *PasswordService) IsValidEmail(email string) bool {

	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
func (p *PasswordService) IsStrongPassword(password string) error {
	var (
		uppercase = `[A-Z]`
		lowercase = `[a-z]`
		number    = `[0-9]`
		special   = `[!@#~$%^&*()_+|<>?:{}]`
	)

	if len(password) < config.MIN_PASSWORD_LENGTH {
		return fmt.Errorf("password must be at least %d characters long", config.MIN_PASSWORD_LENGTH)
	}
	if !regexp.MustCompile(uppercase).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter (A-Z)")
	}
	if !regexp.MustCompile(lowercase).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter (a-z)")
	}
	if !regexp.MustCompile(number).MatchString(password) {
		return errors.New("password must contain at least one number (0-9)")
	}
	if !regexp.MustCompile(special).MatchString(password) {
		return errors.New("password must contain at least one special character (!@#~$%^&*()_+|<>?:{})")
	}

	return nil
}

func (p *PasswordService) Hashpassword(password string) string {

	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashpassword)
}

func (p *PasswordService) ComparePassword(userPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err
}

func (p *PasswordService) PasswordGuidance() string {
	return fmt.Sprintf("Password must be >= %d chars with upper, lower, number, special.", config.MIN_PASSWORD_LENGTH)
}
