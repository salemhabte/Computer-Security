package domain

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	
		
)

type IUserUseCase interface {
	HandleRegistration(user *User) error
	SendOTP(user *User) error
	VerifyOTP(email, otp string) (bool, error)
	Login(email, password string) (*AuthTokens, error)
	Refresh(oldRefreshToken string) (*AuthTokens, error)
	Logout(refreshToken string) error
	UpdateProfile(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	GetUserByEmail(email string) (*UserDTO, error)
}

type IAuthService interface {
	GenerateTokens(user *UserDTO) (string, string, error)
	ValidateRefreshToken(tokenStr string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
	OAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}
type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
}
type IUserOTP interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
}
type IEmailService interface {
	Send(email string, token string) error
	GenerateRandomOTP() string
	SendResetLink(toEmail, subject, message string) error
}

type IUserRepository interface {
	// eka was here
	Create(user *User) error
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CheckUserExistance(userEmail string) bool
	UpdateRole(email, role string) error
	UpdateUserByEmail(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	GetUserByID(userID string) (*UserDTO, error)
	CloseDataBase() error
}
type IAuthRepo interface {
	Save(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	Delete(token string) error
	CloseDataBase() error
}
type IUserOTPRepository interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
	CloseDataBase() error
}

type IOAuthUsecase interface {
	HandleOAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}
type IPasswordUsecase interface {
	GenerateResetToken(email string) error
	ResetPassword(token, newPassword string) error
}