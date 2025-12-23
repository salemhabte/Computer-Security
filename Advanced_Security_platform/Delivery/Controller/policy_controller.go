package controller

import (
	"net/http"
	domain "security/domain"

	"github.com/gin-gonic/gin"
)

type PolicyController struct {
	roleRepo domain.IRoleRepository
	aclRepo  domain.IACLRepository
}

func NewPolicyController(r domain.IRoleRepository, a domain.IACLRepository) *PolicyController {
	return &PolicyController{roleRepo: r, aclRepo: a}
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
	grantor, _ := c.Get("email")
	entry := domain.AccessControlEntry{
		ResourceID: req.ResourceID,
		Subject:    req.Subject,
		Actions:    req.Actions,
		GrantedBy:  grantor.(string),
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
	if err := pc.aclRepo.Revoke(req.ResourceID, req.Subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}

