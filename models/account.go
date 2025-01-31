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
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Username         string             `bson:"username"`
	HashedPassword   string             `bson:"hashed_password"`
	Email            string             `bson:"email"`
	ContactInfo      Contact            `bson:"contact_info"`
	CreatedAt        int64              `bson:"created_at"` // Unix timestamp
	UpdatedAt        int64              `bson:"updated_at"` // Unix timestamp
	Role             Role               `bson:"role"`
	Status           AccountStatus      `bson:"status"`
	RoleSpecificData interface{}        `bson:"role_specific_data"`
	Metadata         Metadata           `bson:"metadata"`
}

type AccountStatus struct {
	IsActive    bool   `bson:"is_active"`
	Deactivated *int64 `bson:"deactivated,omitempty"` // Unix timestamp pointer
	Reason      string `bson:"reason,omitempty"`
}

type Contact struct {
	Phone       string `bson:"phone"`
	Alternative string `bson:"alternative,omitempty"`
	Address     string `bson:"address,omitempty"`
}

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

// Consumer (Your existing structure adapted)
type ConsumerData struct {
	LoyaltyPoints  int             `bson:"loyalty_points"`
	ServiceAddress ServiceAddress  `bson:"service_address"`
	Meters         []Meter         `bson:"meters"`
	BillingProfile BillingProfile  `bson:"billing_profile"`
	PaymentMethods []PaymentMethod `bson:"payment_methods"`
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
