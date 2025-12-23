package conv

import (
	"security/domain"
)

func ChangeToDTOUser(domainUser *domain.User) *domain.UserDTO {
	return &domain.UserDTO{
		UserName:       domainUser.UserName,
		PersonalBio:    domainUser.PersonalBio,
		ProfilePic:     domainUser.ProfilePic,
		Email:          domainUser.Email,
		PhoneNum:       domainUser.PhoneNum,
		TelegramHandle: domainUser.TelegramHandle,
		Password:       domainUser.Password,
		Role:           domainUser.Role,
	}
}
func ChangeToDomainUser(udto *domain.UserDTO) *domain.User {
	return &domain.User{
		UserName:       udto.UserName,
		PersonalBio:    udto.PersonalBio,
		ProfilePic:     udto.ProfilePic,
		Email:          udto.Email,
		PhoneNum:       udto.PhoneNum,
		TelegramHandle: udto.TelegramHandle,
		Password:       udto.Password,
		Role:           udto.Role,
	}
}

func ChangeToDomainVerification(udto *domain.UserUnverifiedDTO) *domain.UserUnverified {
	return &domain.UserUnverified{
		UserName:  udto.UserName,
		Email:     udto.Email,
		OTP:       udto.OTP,
		Password:  udto.Password,
		Role:      udto.Role,
		ExpiresAt: udto.ExpiresAt,
	}
}

func ChangeUnverfiedToVerified(u *domain.UserUnverifiedDTO) *domain.User {
	return &domain.User{
		UserName: u.UserName,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
}
func ChangeToDomainAuthTokens(dto *domain.AuthTokensDTO) *domain.AuthTokens {
	return &domain.AuthTokens{
		AccessToken:  dto.AccessToken,
		RefreshToken: dto.RefreshToken,
	}
}
func ChangeToDomainRefreshToken(dto *domain.RefreshTokenDTO) *domain.RefreshToken {
	return &domain.RefreshToken{
		UserID:    dto.UserID,
		Token:     dto.Token,
		ExpiresAt: dto.ExpiresAt,
	}
}
