package domain

import "time"

// Sensitivity labels for MAC.
type SensitivityLevel string

const (
	SensitivityPublic      SensitivityLevel = "PUBLIC"
	SensitivityInternal    SensitivityLevel = "INTERNAL"
	SensitivityConfidential SensitivityLevel = "CONFIDENTIAL"
)

// Resource describes a protected object.
type Resource struct {
	ID          string
	OwnerEmail  string
	Type        string // e.g. "document", "record"
	Label       SensitivityLevel
	Department  string
	Attributes  map[string]string
}

// Permission describes an action on a resource type.
type Permission struct {
	Name        string // e.g. "document:read"
	Description string
}

// ACL entry for DAC.
type AccessControlEntry struct {
	ResourceID string
	Subject    string   // user email
	Actions    []string // list of permission names
	GrantedBy  string
	GrantedAt  time.Time
}

// Role maps to permissions.
type Role struct {
	Name        string
	Permissions []string
}

// Attribute set used for ABAC.
type AttributeSet struct {
	Department      string
	EmploymentStatus string
	Location        string
	DeviceTrust     string
	Role            string
	BiometricVerified bool
	Custom          map[string]string
}

// Rule conditions for RuBAC/ABAC.
type RuleCondition struct {
	TimeWindow   *TimeWindow
	IPRangeCIDR  []string
	MinDeviceTrust string
	MaxLeaveDays int // example business rule (HR leave approval)
}

type TimeWindow struct {
	Start string // "09:00"
	End   string // "18:00"
	Weekdays []time.Weekday
}

// PolicyDecision describes the result.
type PolicyDecision struct {
	Allowed bool
	Reason  string
}

