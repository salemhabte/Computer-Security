package usecase

import (
	conv "security/Delivery/converter"
	domain "security/domain"
	"net/http"
)

// usecases/user_usecase.go
type OAuthUsecase struct {
	userRepo    domain.IUserRepository
	authService domain.IAuthService
}

func NewOAuthUsecase(ur domain.IUserRepository, as domain.IAuthService) domain.IOAuthUsecase {
	return &OAuthUsecase{
		userRepo:    ur,
		authService: as,
	}
}

func (uc *OAuthUsecase) HandleOAuthLogin(req *http.Request, res http.ResponseWriter) (*domain.UserDTO, error) {
	userData, err := uc.authService.OAuthLogin(req, res)
	if err != nil {
		return nil, err
	}

	existingUser, _ := uc.userRepo.FindByEmail(userData.Email)
	if existingUser != nil {
		return existingUser, nil // Login
	}

	// Signup
	err = uc.userRepo.Create(conv.ChangeToDomainUser(userData))
	if err != nil {
		return nil, err
	}
	return userData, nil
}
