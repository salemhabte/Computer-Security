package usecase

import (
	domain "security/domain"
	"errors"
	"fmt"
	"sync"
	"time"
	"security/config"
)

type UserUsecase struct {
	userinterface domain.IUserRepository
	userVaildate  domain.IUserValidation
	userOTP       domain.IUserOTP
	generateotp   domain.IEmailService
	authService   domain.IAuthService
	authRepo      domain.IAuthRepo
	captcha       domain.ICaptchaValidator

	mu            sync.Mutex
	failedLogins  map[string][]time.Time
	lockedUntil   map[string]time.Time
	mfaPending    map[string]mfaEntry
}

type mfaEntry struct {
	otp     string
	expires time.Time
}

func NewUserUsecase(ui domain.IUserRepository, uv domain.IUserValidation, uo domain.IUserOTP, emailService domain.IEmailService, auth domain.IAuthService, authrepo domain.IAuthRepo, captcha domain.ICaptchaValidator) domain.IUserUseCase {
	return &UserUsecase{
		userinterface: ui,
		userVaildate:  uv,
		userOTP:       uo,
		generateotp:   emailService,
		authService:   auth,
		authRepo:      authrepo,
		captcha:       captcha,
		failedLogins:  make(map[string][]time.Time),
		lockedUntil:   make(map[string]time.Time),
		mfaPending:    make(map[string]mfaEntry),
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

func (a *UserUsecase) Login(email, password, otp, captcha string) (*domain.AuthTokens, error) {
	if !a.captcha.Validate(captcha) {
		return nil, errors.New("captcha validation failed")
	}

	if a.isLocked(email) {
		return nil, errors.New("account locked, try later")
	}

	user, err := a.userinterface.FindByEmail(email)
	if err != nil {
		a.markFail(email)
		return nil, errors.New("user not found")
	}

	err = a.userVaildate.ComparePassword(user.Password, password)
	if err != nil {
		a.markFail(email)
		return nil, errors.New("invalid password")
	}

	// Password ok -> reset fail counter
	a.resetFails(email)

	// MFA flow
	if otp == "" {
		code := a.generateotp.GenerateRandomOTP()
		a.storeMFA(email, code)
		_ = a.generateotp.Send(user.Email, code)
		return nil, errors.New("mfa_required")
	}
	if !a.validateMFA(email, otp) {
		a.markFail(email)
		return nil, errors.New("invalid_or_expired_otp")
	}
	a.clearMFA(email)

	access, refresh, err := a.authService.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

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

func (uc *UserUsecase) GetAttributes(email string) (*domain.AttributeSet, error) {
	return uc.userinterface.GetAttributesByEmail(email)
}

func (a *UserUsecase) markFail(email string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	now := time.Now()
	windowStart := now.Add(-time.Duration(config.LOCKOUT_WINDOW_MIN) * time.Minute)
	a.failedLogins[email] = append(a.failedLogins[email], now)
	var recent []time.Time
	for _, t := range a.failedLogins[email] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}
	a.failedLogins[email] = recent
	if len(recent) >= config.LOCKOUT_THRESHOLD {
		a.lockedUntil[email] = now.Add(time.Duration(config.LOCKOUT_DURATION_MIN) * time.Minute)
	}
}

func (a *UserUsecase) resetFails(email string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.failedLogins[email] = nil
}

func (a *UserUsecase) isLocked(email string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	until, ok := a.lockedUntil[email]
	if !ok {
		return false
	}
	if time.Now().After(until) {
		delete(a.lockedUntil, email)
		return false
	}
	return true
}

func (a *UserUsecase) storeMFA(email, otp string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.mfaPending[email] = mfaEntry{otp: otp, expires: time.Now().Add(5 * time.Minute)}
}

func (a *UserUsecase) validateMFA(email, otp string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	entry, ok := a.mfaPending[email]
	if !ok || time.Now().After(entry.expires) {
		return false
	}
	return entry.otp == otp
}

func (a *UserUsecase) clearMFA(email string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.mfaPending, email)
}