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
	AccessToken  string
	RefreshToken string
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
	UserName       string
	PersonalBio    string
	ProfilePic     string
	Email          string
	PhoneNum       string
	TelegramHandle string
	Password       string
	Role           string
}