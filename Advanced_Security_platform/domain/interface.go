package domain

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type IUserUseCase interface {
	HandleRegistration(user *User) error
	SendOTP(user *User) error
	VerifyOTP(email, otp string) (bool, error)
	Login(email, password, otp, captcha string) (*AuthTokens, error)
	Refresh(oldRefreshToken string) (*AuthTokens, error)
	Logout(refreshToken string) error
	UpdateProfile(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	ChangePassword(email, oldPassword, newPassword string) error
	GetUserByEmail(email string) (*UserDTO, error)
	GetAttributes(email string) (*AttributeSet, error)
}

type IAuthService interface {
	GenerateTokens(user *UserDTO) (string, string, error)
	ValidateRefreshToken(tokenStr string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
	OAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}
type IUserValidation interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) error
	Hashpassword(password string) string
	ComparePassword(userPassword, password string) error
	PasswordGuidance() string
}
type IUserOTP interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
}
type IEmailService interface {
	Send(email string, token string) error
	GenerateRandomOTP() string
	SendResetLink(toEmail, subject, message string) error
}

type IUserRepository interface {
	// eka was here
	Create(user *User) error
	FindByEmail(email string) (*UserDTO, error) //checks if user exisits or not
	UpdatePassword(userID, hashedPassword string) error
	CheckUserExistance(userEmail string) bool
	UpdateRole(email, role string) error
	UpdateUserByEmail(email string, dto *UpdateProfileDTO) (*UserDTO, error)
	GetUserByID(userID string) (*UserDTO, error)
	CloseDataBase() error
	GetAttributesByEmail(email string) (*AttributeSet, error)
}
type IAuthRepo interface {
	Save(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	Delete(token string) error
	CloseDataBase() error
}
type IUserOTPRepository interface {
	StoreOTP(entry UserUnverified) error
	FindOTP(email string) (*UserUnverified, error)
	DeleteOTP(email string) error
	CloseDataBase() error
}

type IOAuthUsecase interface {
	HandleOAuthLogin(req *http.Request, res http.ResponseWriter) (*UserDTO, error)
}
type IPasswordUsecase interface {
	GenerateResetToken(email string) error
	ResetPassword(token, newPassword string) error
}

// Access control and policy
type IRoleRepository interface {
	Upsert(role Role) error
	Get(name string) (*Role, error)
}

type IACLRepository interface {
	Grant(entry AccessControlEntry) error
	Revoke(resourceID, subject string) error
	Get(resourceID, subject string) (*AccessControlEntry, error)
	GetBySubject(subject string) ([]*AccessControlEntry, error)
	GetByResourceID(resourceID string) ([]*AccessControlEntry, error)
}

type IResourceRepository interface {
	GetByID(resourceID string) (*Resource, error)
	Create(resource *Resource) error
	GetByOwner(ownerEmail string) ([]*Resource, error)
	GetByIDs(ids []string) ([]*Resource, error)
}

type IPolicyService interface {
	Decide(user *UserDTO, attrs *AttributeSet, resource *Resource, action string, ctx RequestContext) PolicyDecision
}

type IAuditLogger interface {
	Log(event AuditEvent) error
}

type ICaptchaValidator interface {
	Validate(token string) bool
}

type IBackupService interface {
	RunBackup() error
	LastRun() time.Time
}

// RequestContext carries environmental data for RuBAC/ABAC.
type RequestContext struct {
	IP        string
	UserAgent string
	Now       time.Time
	Location  string
	Device    string
}

type AuditEvent struct {
	UserEmail string                 `json:"user_email"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Result    string                 `json:"result"`
	Reason    string                 `json:"reason"`
	IP        string                 `json:"ip"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}
