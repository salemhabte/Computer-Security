package router

import (
	controller "security/Delivery/Controller"
	"security/infrastructure"

	"github.com/gin-gonic/gin"
)

func Router(uc *controller.UserController, pc *controller.PasswordController, auth *infrastructure.AuthMiddleware,) {
	router := gin.Default()

	router.POST("/login", uc.HandleLogin)
	router.POST("/refresh", uc.HandleRefresh)
	router.POST("/registration", uc.Registration)
	router.POST("/registration/verification", uc.RegistrationValidation)
	router.POST("/forgot_password", pc.ForgotPassword)
	router.PUT("/reset_password", pc.ResetPassword)

	userRoutes := router.Group("/user")
	userRoutes.Use(auth.JWTAuthMiddleware())
	{
		userRoutes.POST("/logout", uc.HandleLogout)
		userRoutes.PUT("/edit_profile", uc.UpdateProfile)
	}

	router.GET("/auth/:provider", uc.SignInWithProvider)
	router.GET("/auth/:provider/callback", uc.CallbackHandler)
	router.GET("/success", uc.Success)

	router.Run()
}