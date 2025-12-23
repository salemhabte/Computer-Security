package usecase

import (
	domain "security/domain"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type passwordUsecase struct {
	userRepo  		domain.IUserRepository
	emailService	domain.IEmailService
	jwtSecret 		string
}

func NewPasswordUsecase(repo domain.IUserRepository, emailser domain.IEmailService, jwtSec string) domain.IPasswordUsecase {
	return &passwordUsecase{
		userRepo:  repo,
		emailService: emailser,
		jwtSecret: jwtSec,
	}
}

func (u *passwordUsecase) GenerateResetToken(email string) error {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	claims := jwt.MapClaims{
		"user_id": user.UserID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))

	if err != nil {
		return err
	}

	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", signedToken)
	subject := "Password Reset Request"
	message := fmt.Sprintf("Click the link below to reset your password:\n\n%s\n\nThis link will expire in 15 minutes.", resetURL)
	
	if err := u.emailService.SendResetLink(user.Email, subject, message); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}
	return nil
}

func (u *passwordUsecase) ResetPassword(tokenStr, newPassword string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("invalid email in token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return u.userRepo.UpdatePassword(email, string(hashedPassword))

}
