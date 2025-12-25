package controller

import (
	"net/http"
	domain "security/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (pc *PolicyController) CreateResource(c *gin.Context) {
	var req struct {
		Type       string            `json:"type"`
		Label      string            `json:"label"` // PUBLIC, INTERNAL, CONFIDENTIAL
		Attributes map[string]string `json:"attributes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail, _ := c.Get("email")
	userDept, _ := c.Get("department") // Assuming department is in context claims or we fetch it?
	// The JWT claims in `userUsecase` might not put Department.
	// `AuthMiddleware` extracts claims.
	// Let's assume we can get it or just lookup user.
	// For MVP, if "department" isn't in claims, we might need to fetch user.
	// But let's check `User` struct in `domain`.
	// `domain.User` has `Department`.
	// If `userDept` is nil, we default or leave empty.

	// Better: just use what we have.

	entry := &domain.Resource{
		ID:         primitive.NewObjectID().Hex(),
		OwnerEmail: userEmail.(string),
		Type:       req.Type,
		Label:      domain.SensitivityLevel(req.Label),
		Department: "", // We might need to fetch user profile to fill this. For now leaving empty or "General".
		Attributes: req.Attributes,
	}

	if val, ok := userDept.(string); ok {
		entry.Department = val
	}

	if err := pc.resourceRepo.Create(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "resource created", "id": entry.ID})
}

func (pc *PolicyController) ListMyResources(c *gin.Context) {
	userEmail, _ := c.Get("email")
	email := userEmail.(string)

	// 1. Get Owned Resources
	ownedResources, err := pc.resourceRepo.GetByOwner(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Get Shared Resources via ACL
	acls, err := pc.aclRepo.GetBySubject(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sharedIDs []string
	for _, acl := range acls {
		sharedIDs = append(sharedIDs, acl.ResourceID)
	}

	// Fetch details for shared resources
	sharedResources, err := pc.resourceRepo.GetByIDs(sharedIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Combine and Deduplicate
	resourceMap := make(map[string]*domain.Resource)
	for _, r := range ownedResources {
		resourceMap[r.ID] = r
	}
	for _, r := range sharedResources {
		// Only add if not already present (owned)
		if _, exists := resourceMap[r.ID]; !exists {
			resourceMap[r.ID] = r
		}
	}

	var allResources []*domain.Resource
	for _, r := range resourceMap {
		allResources = append(allResources, r)
	}

	c.JSON(http.StatusOK, allResources)
}

func (pc *PolicyController) GetResourcePermissions(c *gin.Context) {
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource id required"})
		return
	}

	userEmail, _ := c.Get("email")
	userRole, _ := c.Get("role")
	email := userEmail.(string)
	role := userRole.(string)

	// Security Check: Only Owner or Admin can list permissions
	// Be strict as per latest user request
	if pc.resourceRepo != nil {
		resource, err := pc.resourceRepo.GetByID(resourceID)
		if err != nil || resource == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		if resource.OwnerEmail != email && role != "SUPER_ADMIN" && role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only owner can view permissions"})
			return
		}
	}

	acls, err := pc.aclRepo.GetByResourceID(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acls)
}
