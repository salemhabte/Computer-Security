package domain

import (
	"time"
)

const (
	ADMIN       = "ADMIN"
	UESR        = "USER"
	SUPER_ADMIN = "SUPER_ADMIN"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshToken struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}
type UserUnverified struct {
	UserName  string
	Email     string
	OTP       string
	Password  string
	Role      string
	ExpiresAt time.Time
}

type User struct {
	UserName          string
	PersonalBio       string
	ProfilePic        string
	Email             string
	PhoneNum          string
	TelegramHandle    string
	Password          string
	Role              string
	Department        string
	EmploymentStatus  string
	DeviceTrust       string
	BiometricVerified bool
}
