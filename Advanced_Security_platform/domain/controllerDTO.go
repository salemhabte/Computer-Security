package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenDTO struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type UserUnverifiedDTO struct {
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	Password  string    `json:"password" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type UserDTO struct {
	UserID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName       string             `json:"username" bson:"username"`
	PersonalBio    string             `json:"personal_bio" bson:"personal_bio"`
	ProfilePic     string             `json:"profile_pic" bson:"profile_pic"` // store as URL or base64
	Email          string             `json:"email" bson:"email"`
	PhoneNum       string             `json:"phone_num" bson:"phone_num"` // validate format
	TelegramHandle string             `json:"telegram_handle" bson:"telegram_handle"`
	Password       string             `json:"password" bson:"password"` // securely hashed
	Role           string             `json:"role" bson:"role"`         // values: "admin", "user", "super_admin"
}


type UpdateProfileDTO struct {
	UserName       string             `json:"username" bson:"username"`
	PersonalBio    string             `json:"personal_bio" bson:"personal_bio"`
	ProfilePic     string             `json:"profile_pic" bson:"profile_pic"` // store as URL or base64
	Email          string             `json:"email" bson:"email"`
	PhoneNum       string             `json:"phone_num" bson:"phone_num"` // validate format
	TelegramHandle string             `json:"telegram_handle" bson:"telegram_handle"`
	Password       string             `json:"password" bson:"password"` // securely hashed
}

