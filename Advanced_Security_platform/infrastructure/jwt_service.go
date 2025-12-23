package infrastructure

import (
	domain "security/domain"
	"security/config"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth/gothic"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService struct {
	authRepo domain.IAuthRepo
}

func NewJWTService(auth domain.IAuthRepo) *JWTService {
	return &JWTService{
		authRepo: auth,
	}
}

var jwtSecret = []byte(config.JWTSECRET)
var refreshSecret = []byte(config.JWTREFRESHSECRET)

func (j *JWTService) GenerateTokens(user *domain.UserDTO) (string, string, error) {
	if user.Email == "" {
		return "", "", errors.New("user email cannot be empty")
	}
	// Access Token
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	atString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	rtClaims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return atString, rtString, nil
}

func (j *JWTService) ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or malformed token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Manually check expiration
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("invalid exp claim")
	}
	if int64(expFloat) < time.Now().Unix() {
		_ = j.authRepo.Delete(tokenStr)
		return "", errors.New("refresh token expired")
	}

	// Extract user_id safely
	userIDRaw, ok := claims["user_id"]
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		return "", errors.New("user_id is not a string")
	}

	return userID, nil
}

func (j *JWTService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return nil, errors.New("invalid claims")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("invalid or missing email in token")
	}

	return claims, nil
}

func (o *JWTService) OAuthLogin(req *http.Request, res http.ResponseWriter) (*domain.UserDTO, error) {
	user, err := gothic.CompleteUserAuth(res, req)
	fmt.Print("+++++++", req, "+++++++", err, "------")
	if err != nil {
		return nil, err
	}

	return &domain.UserDTO{
		Email:    user.Email,
		UserName: user.Name,
	}, nil
}
