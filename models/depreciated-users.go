package models

import (
	"time"
)

// Base User struct that contains common fields for all user types
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string
}

// Define structs for each role that contain role-specific fields
type SystemAdmin struct {
	User
	AdminLevel int
}

type FinancialAdmin struct {
	User
	BudgetAccess bool
}

type HRAdmin struct {
	User
	EmployeeAccess bool
}

type CustomerServiceAdmin struct {
	User
	SupportLevel int
}

type OperationsMonitor struct {
	User
	MonitorLevel int
}

type Cashier struct {
	User
	TillNumber int
}

type Consumer struct {
	User
	LoyaltyPoints int
}

// Define constructor functions for each role.
func NewSystemAdmin(id int, username, password, email string) *SystemAdmin {
	return &SystemAdmin{
		User: User{
			ID:        id,
			Username:  username,
			Password:  password,
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Role:      "SystemAdmin",
		},
		AdminLevel: 1,
	}
}

func NewFinancialAdmin(id int, username, password, email string, budgetAccess bool) *FinancialAdmin {
	return &FinancialAdmin{
		User: User{
			ID:        id,
			Username:  username,
			Password:  password,
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Role:      "FinancialAdmin",
		},
		BudgetAccess: budgetAccess,
	}
}
