package controller

import (
	"fmt"
	"net/http"
	"security/domain"

	"security/Delivery/converter"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)
type UserController struct {
	UserUsecase  domain.IUserUseCase
	OauthUsecase domain.IOAuthUsecase
}
func NewUserController(uuc domain.IUserUseCase, oat domain.IOAuthUsecase) *UserController {
	return &UserController{
		UserUsecase:  uuc,
		OauthUsecase: oat,
	}
}

func (uc *UserController) Registration(ctx *gin.Context) {

	var user *domain.UserDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	
	err := uc.UserUsecase.HandleRegistration(conv.ChangeToDomainUser(user))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Please enter your otp to successfully register",
	})
}

func (uc *UserController) RegistrationValidation(ctx *gin.Context) {
	var user *domain.UserUnverifiedDTO

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if user.Email == "" || user.OTP == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	userone := conv.ChangeToDomainVerification(user)
	valid, err := uc.UserUsecase.VerifyOTP(userone.Email, userone.OTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while verifying OTP",
		})
		return
	}

	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired OTP",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

func (uc *UserController) HandleLogin(ctx *gin.Context) {

	var req struct {
		Email   string `json:"email"`
		Password string `json:"password"`
		OTP     string `json:"otp"`
		Captcha string `json:"captcha_token"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	jwtToken, err := uc.UserUsecase.Login(req.Email, req.Password, req.OTP, req.Captcha)
	if err != nil {
		if err.Error() == "mfa_required" {
			ctx.JSON(http.StatusAccepted, gin.H{"message": "OTP sent to email. Provide otp to complete login."})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User logged in successfully", "token": jwtToken})

}
func (h *UserController) HandleRefresh(c *gin.Context) {
	var input *domain.RefreshTokenDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	tokens, err := h.UserUsecase.Refresh(input.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens) // new AuthTokensDTO
}

func (h *UserController) HandleLogout(c *gin.Context) {
	var input *domain.RefreshTokenDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	err := h.UserUsecase.Logout(input.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}


func (uc *UserController) SignInWithProvider(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	// req := c.Request
	// req = req.WithContext(context.WithValue(req.Context(), "provider", provider))
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (uc *UserController) CallbackHandler(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	// req := c.Request
	// fmt.Println("^^^^^",provider)
	// req = req.WithContext(context.WithValue(c.Request.Context(), "provider", provider))

	user, err := uc.OauthUsecase.HandleOAuthLogin(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//   _, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Logged in", "user": user})

	// user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	// if err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	c.Redirect(http.StatusTemporaryRedirect, "/success")
}
func (uc *UserController) Success(c *gin.Context) {

	c.Data(http.StatusOK, "text/html; charset=utf-8", fmt.Appendf(nil, `
      <div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfull signed in!</h1>
          
          </div>
      </div>
  `))
}

func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	var updateDTO domain.UpdateProfileDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	emailVal, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email in context"})
		return
	}
	updatedUser, err := uc.UserUsecase.UpdateProfile(email, &updateDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated", "user": updatedUser})
}
