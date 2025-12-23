package infrastructure

import (
	domain "security/domain"
	"time"

	"github.com/gin-gonic/gin"
)

func PolicyMiddleware(policy domain.IPolicyService, audit domain.IAuditLogger, resProvider func(c *gin.Context) *domain.Resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		emailVal, _ := c.Get("email")
		roleVal, _ := c.Get("role")
		userID, _ := c.Get("user_id")

		resource := resProvider(c)
		email, _ := emailVal.(string)
		role, _ := roleVal.(string)
		user := &domain.UserDTO{
			Email: email,
			Role:  role,
		}
		ctx := domain.RequestContext{
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Now:       time.Now(),
		}

		// attrs may be fetched lazily by caller; here nil implies minimal info
		decision := policy.Decide(user, &domain.AttributeSet{Role: user.Role}, resource, action, ctx)
		if audit != nil {
			_ = audit.Log(domain.AuditEvent{
				UserEmail: user.Email,
				Action:    action,
				Resource:  func() string {
					if resource != nil {
						return resource.ID
					}
					return ""
				}(),
				Result:    map[bool]string{true: "allow", false: "deny"}[decision.Allowed],
				Reason:    decision.Reason,
				IP:        ctx.IP,
				Timestamp: ctx.Now,
			})
		}

		if !decision.Allowed {
			c.JSON(403, gin.H{"error": decision.Reason})
			c.Abort()
			return
		}
		c.Next()
		_ = userID
	}
}

