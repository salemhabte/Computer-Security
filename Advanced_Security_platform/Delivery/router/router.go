package router

import (
	controller "security/Delivery/Controller"
	domain "security/domain"
	"security/infrastructure"

	"github.com/gin-gonic/gin"
)

func Router(uc *controller.UserController, pc *controller.PasswordController, polc *controller.PolicyController, bc *controller.BackupController, auth *infrastructure.AuthMiddleware, policy domain.IPolicyService, audit domain.IAuditLogger) {
	router := gin.Default()
	router.Use(CORSMiddleware())

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
		userRoutes.POST("/change_password", uc.HandleChangePassword)
	}

	// Admin / policy routes
	admin := router.Group("/admin")
	admin.Use(auth.JWTAuthMiddleware(), infrastructure.RoleMiddleware("SUPER_ADMIN"))
	{
		admin.POST("/role", polc.UpsertRole)
		admin.POST("/backup/run", bc.RunBackup)
	}

	// DAC routes - owners OR admins can grant/revoke
	dacRoutes := router.Group("/dac")
	dacRoutes.Use(auth.JWTAuthMiddleware())
	{
		dacRoutes.POST("/grant", polc.GrantDAC)
		dacRoutes.POST("/revoke", polc.RevokeDAC)
		dacRoutes.POST("/resources", polc.CreateResource)
		dacRoutes.GET("/resources", polc.ListMyResources)
		dacRoutes.GET("/permissions/:id", polc.GetResourcePermissions)
	}

	router.GET("/auth/:provider", uc.SignInWithProvider)
	router.GET("/auth/:provider/callback", uc.CallbackHandler)
	router.GET("/success", uc.Success)

	router.Run()
}
