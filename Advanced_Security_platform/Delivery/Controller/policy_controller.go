package controller

import (
	"net/http"
	domain "security/domain"

	"github.com/gin-gonic/gin"
)

type PolicyController struct {
	roleRepo     domain.IRoleRepository
	aclRepo      domain.IACLRepository
	resourceRepo domain.IResourceRepository
}

func NewPolicyController(r domain.IRoleRepository, a domain.IACLRepository, resRepo domain.IResourceRepository) *PolicyController {
	return &PolicyController{roleRepo: r, aclRepo: a, resourceRepo: resRepo}
}

func (pc *PolicyController) UpsertRole(c *gin.Context) {
	var role domain.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if role.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name required"})
		return
	}
	if err := pc.roleRepo.Upsert(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role upserted"})
}

func (pc *PolicyController) GrantDAC(c *gin.Context) {
	var req struct {
		ResourceID string   `json:"resource_id"`
		Subject    string   `json:"subject"`
		Actions    []string `json:"actions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grantorEmail, _ := c.Get("email")
	grantorRole, _ := c.Get("role")
	grantor := grantorEmail.(string)
	
	// DAC: Check if user is resource owner OR admin
	if pc.resourceRepo != nil {
		resource, err := pc.resourceRepo.GetByID(req.ResourceID)
		if err == nil && resource != nil {
			// Check ownership
			if resource.OwnerEmail != grantor {
				// Not owner, check if admin
				if grantorRole != "SUPER_ADMIN" && grantorRole != "ADMIN" {
					c.JSON(http.StatusForbidden, gin.H{"error": "only resource owner or admin can grant permissions"})
					return
				}
			}
		}
	} else {
		// If no resource repo, fallback to admin-only
		if grantorRole != "SUPER_ADMIN" && grantorRole != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
	}
	
	entry := domain.AccessControlEntry{
		ResourceID: req.ResourceID,
		Subject:    req.Subject,
		Actions:    req.Actions,
		GrantedBy:  grantor,
	}
	if err := pc.aclRepo.Grant(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "granted"})
}

func (pc *PolicyController) RevokeDAC(c *gin.Context) {
	var req struct {
		ResourceID string `json:"resource_id"`
		Subject    string `json:"subject"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grantorEmail, _ := c.Get("email")
	grantorRole, _ := c.Get("role")
	grantor := grantorEmail.(string)
	
	// DAC: Check if user is resource owner OR admin
	if pc.resourceRepo != nil {
		resource, err := pc.resourceRepo.GetByID(req.ResourceID)
		if err == nil && resource != nil {
			// Check ownership
			if resource.OwnerEmail != grantor {
				// Not owner, check if admin
				if grantorRole != "SUPER_ADMIN" && grantorRole != "ADMIN" {
					c.JSON(http.StatusForbidden, gin.H{"error": "only resource owner or admin can revoke permissions"})
					return
				}
			}
		}
	} else {
		// If no resource repo, fallback to admin-only
		if grantorRole != "SUPER_ADMIN" && grantorRole != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
	}
	
	if err := pc.aclRepo.Revoke(req.ResourceID, req.Subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}

