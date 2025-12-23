package infrastructure

import (
	domain "security/domain"
	"strings"
	"time"
)

// PolicyEngine implements MAC + DAC + RBAC + ABAC + RuBAC.
type PolicyEngine struct {
	roleRepo domain.IRoleRepository
	aclRepo  domain.IACLRepository
}

func NewPolicyEngine(roleRepo domain.IRoleRepository, aclRepo domain.IACLRepository) domain.IPolicyService {
	return &PolicyEngine{
		roleRepo: roleRepo,
		aclRepo:  aclRepo,
	}
}

func (p *PolicyEngine) Decide(user *domain.UserDTO, attrs *domain.AttributeSet, res *domain.Resource, action string, ctx domain.RequestContext) domain.PolicyDecision {
	// MAC: enforce sensitivity vs role
	if res != nil && !p.macCheck(user, res) {
		return domain.PolicyDecision{Allowed: false, Reason: "MAC: insufficient clearance"}
	}

	// RuBAC: time/device/IP checks
	if !p.rubacCheck(ctx) {
		return domain.PolicyDecision{Allowed: false, Reason: "RuBAC: outside allowed time or context"}
	}

	// DAC: owner or ACL entry
	if res != nil && res.OwnerEmail == user.Email {
		return domain.PolicyDecision{Allowed: true, Reason: "DAC: owner access"}
	}
	if res != nil {
		if ace, err := p.aclRepo.Get(res.ID, user.Email); err == nil && ace != nil {
			for _, act := range ace.Actions {
				if act == action {
					return domain.PolicyDecision{Allowed: true, Reason: "DAC: ACL permit"}
				}
			}
		}
	}

	// RBAC: role permission mapping
	if user != nil && user.Role != "" {
		role, err := p.roleRepo.Get(user.Role)
		if err == nil && role != nil {
			for _, perm := range role.Permissions {
				if perm == action {
					return domain.PolicyDecision{Allowed: true, Reason: "RBAC: role permit"}
				}
			}
		}
	}

	// ABAC: department/status/device trust
	if !p.abacCheck(attrs, res, action, ctx) {
		return domain.PolicyDecision{Allowed: false, Reason: "ABAC: attribute constraints"}
	}

	// Default deny
	return domain.PolicyDecision{Allowed: false, Reason: "default deny"}
}

func (p *PolicyEngine) macCheck(user *domain.UserDTO, res *domain.Resource) bool {
	if res == nil || user == nil {
		return true
	}
	// Simple clearance: SUPER_ADMIN can access all, ADMIN all except confidential write, USER no confidential
	if strings.EqualFold(user.Role, domain.SUPER_ADMIN) {
		return true
	}
	if strings.EqualFold(user.Role, domain.ADMIN) {
		return res.Label != domain.SensitivityConfidential
	}
	return res.Label == domain.SensitivityPublic || res.Label == domain.SensitivityInternal
}

func (p *PolicyEngine) rubacCheck(ctx domain.RequestContext) bool {
	// Basic time window: allow only between 06:00 and 22:00 local
	now := ctx.Now
	if now.IsZero() {
		now = time.Now()
	}
	hour := now.Hour()
	return hour >= 6 && hour <= 22
}

func (p *PolicyEngine) abacCheck(attrs *domain.AttributeSet, res *domain.Resource, action string, ctx domain.RequestContext) bool {
	if attrs == nil {
		return false
	}
	// Example: payroll data only for payroll dept; finance managers during business hours
	if res != nil && strings.EqualFold(res.Department, "Payroll") && attrs.Department != "Payroll" {
		return false
	}
	if strings.Contains(action, "approve_leave") && attrs.Role != "HR_MANAGER" {
		return false
	}
	// Device trust minimal check
	if attrs.DeviceTrust != "" && attrs.DeviceTrust == "LOW" {
		return false
	}
	// Require biometric for high-sensitivity actions
	if res != nil && res.Label == domain.SensitivityConfidential && !attrs.BiometricVerified {
		return false
	}
	_ = ctx
	return true
}

