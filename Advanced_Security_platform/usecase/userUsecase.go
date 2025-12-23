package usecase

import (
	domain "security/domain"
	"errors"
	"fmt"
	"time"
)

type UserUsecase struct {
	userinterface domain.IUserRepository
	userVaildate  domain.IUserValidation
	userOTP       domain.IUserOTP
	generateotp   domain.IEmailService
	authService   domain.IAuthService
	authRepo      domain.IAuthRepo
}

func NewUserUsecase(ui domain.IUserRepository, uv domain.IUserValidation, uo domain.IUserOTP, emailService domain.IEmailService, auth domain.IAuthService, authrepo domain.IAuthRepo) domain.IUserUseCase {
	return &UserUsecase{
		userinterface: ui,
		userVaildate:  uv,
		userOTP:       uo,
		generateotp:   emailService,
		authService:   auth,
		authRepo:      authrepo,
	}
}

func (uc *UserUsecase) HandleRegistration(user *domain.User) error {
	existing := uc.userinterface.CheckUserExistance(user.Email)

	if existing {
		return errors.New("user already exists")
	}

	isvaild_email := uc.userVaildate.IsValidEmail(user.Email)
	ispassword_strong := uc.userVaildate.IsStrongPassword(user.Password)

	if !ispassword_strong || !isvaild_email {
		return errors.New("invalid password or email")
	}

	hashpassword := uc.userVaildate.Hashpassword(user.Password)
	user.Password = hashpassword

	err := uc.SendOTP(user)

	if err != nil {
		return errors.New("failed to send OTP: " + err.Error())
	}

	return nil
}
func (uc *UserUsecase) SendOTP(user *domain.User) error {
	otp := uc.generateotp.GenerateRandomOTP()
	entry := domain.UserUnverified{
		UserName:  user.UserName,
		Email:     user.Email,
		OTP:       otp,
		Password:  user.Password,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := uc.generateotp.Send(user.Email, otp); err != nil {
		return err
	}
	return uc.userOTP.StoreOTP(entry)
}

func (uc *UserUsecase) VerifyOTP(email, otp string) (bool, error) {
	entry, err := uc.userOTP.FindOTP(email)
	if err != nil || entry == nil {
		return false, err
	}
	if time.Now().After(entry.ExpiresAt) || entry.OTP != otp {
		return false, nil
	}
	verifiedUser := &domain.User{
		UserName: entry.UserName,
		Email:    entry.Email,
		Password: entry.Password,
		Role:     entry.Role,
	}
	err = uc.userinterface.Create(verifiedUser)

	if err != nil {
		return false, err
	}
	_ = uc.userOTP.DeleteOTP(email)
	return true, nil
}

func (uc *UserUsecase) PromoteUser(actor, target string) error {
	actUser, err := uc.userinterface.FindByEmail(actor)

	if err != nil {
		return fmt.Errorf("user not found")
	}

	if actUser.Role != "SUPER_ADMIN" {
		return fmt.Errorf("unauthorized user")
	}

	return uc.userinterface.UpdateRole(target, "ADMIN")
}

func (uc *UserUsecase) DemoteUser(actor, target string) error {
	actUser, err := uc.userinterface.FindByEmail(actor)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if actUser.Role != "SUPER_ADMIN" {
		return fmt.Errorf("unauthorized user")
	}

	return uc.userinterface.UpdateRole(target, "USER")
}

func (a *UserUsecase) Login(email, password string) (*domain.AuthTokens, error) {
	user, err := a.userinterface.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = a.userVaildate.ComparePassword(user.Password, password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	access, refresh, err := a.authService.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	// Save refresh token in DB
	refreshEntry := &domain.RefreshToken{
		UserID:    user.UserID.Hex(),
		Token:     refresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = a.authRepo.Save(refreshEntry)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}
func (a *UserUsecase) Refresh(oldRefreshToken string) (*domain.AuthTokens, error) {
	// Validate token structure
	userID, err := a.authService.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return nil, err
	}

	// Check if token is stored in DB
	stored, err := a.authRepo.GetByToken(oldRefreshToken)
	if err != nil || stored.UserID != userID {
		return nil, errors.New("refresh token not found or mismatched")
	}

	// Optional: delete old token (rotation)
	_ = a.authRepo.Delete(oldRefreshToken)

	user, err := a.userinterface.GetUserByID(userID)

	if err != nil {
		return nil, errors.New("user not found by id")
	}

	// Generate new tokens
	newAccess, newRefresh, err := a.authService.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	newEntry := &domain.RefreshToken{
		UserID:    userID,
		Token:     newRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err = a.authRepo.Save(newEntry)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokens{AccessToken: newAccess, RefreshToken: newRefresh}, nil
}

func (uc *UserUsecase) Logout(refreshToken string) error {
	return uc.authRepo.Delete(refreshToken)
}

func (uc *UserUsecase) UpdateProfile(email string, dto *domain.UpdateProfileDTO) (*domain.UserDTO, error) {
	return uc.userinterface.UpdateUserByEmail(email, dto)
}

// Just newly added
func (uc *UserUsecase) GetUserByEmail(email string) (*domain.UserDTO, error) {
	return uc.userinterface.FindByEmail(email)
}
