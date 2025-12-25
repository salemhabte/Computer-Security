package usecase

import "errors"

func (uc *UserUsecase) ChangePassword(email, oldPassword, newPassword string) error {
	user, err := uc.userinterface.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// 1. Verify old password
	// FindByEmail returns UserDTO which might not have password hash?
	// DTO usually hides password. Let's check FindByEmail implementation or UserDTO struct.
	// Actually user.Password in DTO might be hidden.
	// I need the stored hash.
	// domain.User has Password. domain.UserDTO usually doesn't or it's empty.
	// Let's check domain.go for UserDTO.

	// Assuming FindByEmail returns DTO, I might need to fetch the full User struct or
	// if DTO has it.
	// Wait, standard practice: DTO shouldn't have it.
	// I might need `userinterface.GetUserByID` or similar that returns full User or specific method.
	// However, `FindByEmail` sig is (*UserDTO, error).
	// Let's assume for a second I can't get password from DTO.
	// But `Login` uses `uc.userinterface.FindByEmail(email)`.
	// Let's check `Login` implementation in `userUsecase.go`.
	// Line 147: `user, err := a.userinterface.FindByEmail(email)`
	// Line 153: `err = a.userVaildate.ComparePassword(user.Password, password)`
	// So `UserDTO` DOES have the password hash.

	if err := uc.userVaildate.ComparePassword(user.Password, oldPassword); err != nil {
		return errors.New("incorrect old password")
	}

	// 2. Validate new password
	if err := uc.userVaildate.IsStrongPassword(newPassword); err != nil {
		return err
	}

	// 3. Hash new password
	hashed := uc.userVaildate.Hashpassword(newPassword)

	// 4. Update
	// user.UserID is string or ObjectID?
	// In Login it uses user.UserID.Hex() for token.
	// UpdatePassword expects (userID, hashed string).
	return uc.userinterface.UpdatePassword(user.UserID.Hex(), hashed)
}
