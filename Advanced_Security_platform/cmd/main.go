package main

import (
	"log"
	"security/Delivery/oauth"
	"security/Delivery/Controller"
	"security/Delivery/router"
	repositories "security/Repository"
	"security/config"
	"security/infrastructure"
	usecases"security/usecase"
)

func main() {
	config.InitEnv()
	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("❌Failed to connect:", err)
	}
	defer mongoClient.Disconnect()

	from := config.FROM
	appPass := config.APPPASS
	smtpServer := config.SMTPSERVER
	smtpPort := config.SMTPPORT
	user := config.SMTPUSER
	jwtSecret := config.JWTSECRET

	oauth.InitOAuth()
	// Initialize DataBase Repository
	userRepo := repositories.NewUserRepository()
	authRepo := repositories.NewRefreshTokenRepository()
	otpService := repositories.NewUserOTPRepository()
	roleRepo := repositories.NewRoleRepository()
	aclRepo := repositories.NewACLRepository()
	resourceRepo := repositories.NewResourceRepository()

	// Initialize Services
	passwaordService := infrastructure.NewPasswordService()
	captchaValidator := infrastructure.NewCaptchaValidator()
	authMiddleware := infrastructure.NewAuthMiddleware(authRepo)
	authService := infrastructure.NewJWTService(authRepo)
	emailService := infrastructure.NewOTP_service(from, appPass, smtpServer, smtpPort, user)
	policyEngine := infrastructure.NewPolicyEngine(roleRepo, aclRepo)
	auditLogger, _ := infrastructure.NewAuditLogger("audit.log")
	backupService := infrastructure.NewBackupService()

	// Initialize UseCases
	oauthUsecase := usecases.NewOAuthUsecase(userRepo, authService)
	userUsecase := usecases.NewUserUsecase(userRepo, passwaordService, otpService, emailService, authService, authRepo, captchaValidator)
	passwordUsecase := usecases.NewPasswordUsecase(userRepo, emailService, jwtSecret)

	// Initialize DataBase Repository
	userController := controller.NewUserController(userUsecase, oauthUsecase)
	passwordController := controller.NewPasswordController(passwordUsecase)
	policyController := controller.NewPolicyController(roleRepo, aclRepo, resourceRepo)
	backupController := controller.NewBackupController(backupService)
	



	// Initialize Routers
	router.Router(userController, passwordController, policyController, backupController, authMiddleware, policyEngine, auditLogger,)

}
