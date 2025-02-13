package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---Unified Base Account Structure---//
type Role string

const (
	RoleSystemAdmin          Role = "system_admin"
	RoleFinancialAdmin       Role = "financial_admin"
	RoleHRAdmin              Role = "hr_admin"
	RoleCustomerServiceAdmin Role = "customer_service_admin"
	RoleOperationsMonitor    Role = "operations_monitor"
	RoleCashier              Role = "cashier"
	RoleConsumer             Role = "consumer"
)

type Account struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" validate:"-"`
	HashedPassword     string             `bson:"hashed_password" validate:"required,min=8"`                                                                              // Password must be at least 8 characters
	Email              string             `bson:"email" validate:"required,email"`                                                                                        // Email must be valid
	CreatedAt          int64              `bson:"created_at" validate:"required"`                                                                                         // Must be a valid timestamp
	UpdatedAt          int64              `bson:"updated_at" validate:"required"`                                                                                         // Must be a valid timestamp
	Role               Role               `bson:"role" validate:"oneof=system_admin financial_admin hr_admin customer_service_admin operations_monitor cashier consumer"` // Only predefined roles allowed
	Status             AccountStatus      `bson:"status" validate:"required"`
	RoleSpecificDataID primitive.ObjectID `bson:"role_specific_data_id" validate:"required"` // Must be a valid ObjectID
}

type AccountStatus struct {
	IsActive    bool   `bson:"is_active" validate:"required"`                 // Must be true or false
	Deactivated *int64 `bson:"deactivated,omitempty" validate:"omitempty"`    // Optional, but must be a valid timestamp if present
	Reason      string `bson:"reason,omitempty" validate:"omitempty,max=255"` // Optional, but max length 255 if present
}

/*

type Metadata struct {
	LastLogin   int64             `bson:"last_login"`
	Preferences map[string]string `bson:"preferences"`
	AuditTrail  []AuditEvent      `bson:"audit_trail"`
	SystemNotes []string          `bson:"system_notes"`
	APIKeys     []APIKey          `bson:"api_keys"`
}

// ---Role-Specific Structures---//
// System Admin
type SystemAdminData struct {
	AccessLevel       int      `bson:"access_level"`
	AssignedRegions   []string `bson:"assigned_regions"`
	SecurityClearance int      `bson:"security_clearance"`
	MFAEnforced       bool     `bson:"mfa_enforced"`
}

// Financial Admin
type FinancialAdminData struct {
	BudgetAreas     []string `bson:"budget_areas"`
	ApprovalLimit   float64  `bson:"approval_limit"`
	AccessLevel     string   `bson:"access_level"`
	AuditPrivileges bool     `bson:"audit_privileges"`
}

// Cashier
type CashierData struct {
	TillNumber     int     `bson:"till_number"`
	ShiftSchedule  string  `bson:"shift_schedule"`
	TransactionCap float64 `bson:"transaction_cap"`
}

// Operations Monitor
type OperationsMonitorData struct {
	MonitoredSystems []string `bson:"monitored_systems"`
	AlertLevel       int      `bson:"alert_level"`
	AccessFeeds      []string `bson:"access_feeds"`
}

// ---Supporting Structures---//
type AuditEvent struct {
	Timestamp int64  `bson:"timestamp"`
	Actor     string `bson:"actor"`
	Action    string `bson:"action"`
	IPAddress string `bson:"ip_address"`
}

type APIKey struct {
	Key         string   `bson:"key"`
	CreatedAt   int64    `bson:"created_at"`
	ExpiresAt   int64    `bson:"expires_at"`
	Permissions []string `bson:"permissions"`
}

type ServiceAddress struct {
	Location    GeoJSON `bson:"location"`
	PostalCode  string  `bson:"postal_code"`
	ServiceType string  `bson:"service_type"`
}

*/
