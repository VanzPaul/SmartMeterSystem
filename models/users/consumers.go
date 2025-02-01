package users

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Enums for type safety
type AccountType string

const (
	AccountTypeConsumer AccountType = "consumer"
	AccountTypeBusiness AccountType = "business"
)

type ConsumerType string

const (
	ConsumerTypeResidential ConsumerType = "residential"
	ConsumerTypeCommercial  ConsumerType = "commercial"
)

type CommandType string

const (
	CommandTypeMeterRead CommandType = "METER_READ"
	CommandTypeReboot    CommandType = "REBOOT"
)

type CommandStatus string

const (
	CommandStatusPending   CommandStatus = "PENDING"
	CommandStatusCompleted CommandStatus = "COMPLETED"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// Core account structure
type Consumer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AccountNumber string             `bson:"accountNumber"`
	Name          string             `bson:"name"`
	Address       Address            `bson:"address"`
	Contact       Contact            `bson:"contact"`
	AccountType   AccountType        `bson:"accountType"`
	ConsumerType  ConsumerType       `bson:"consumerType"`
	Meters        []Meter            `bson:"meters"`
	Billing       Billing            `bson:"billing"`
	Ledger        Ledger             `bson:"ledger"`
	CreatedAt     int64              `bson:"createdAt"`
	UpdatedAt     int64              `bson:"updatedAt"`
}

type Address struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	State      string `bson:"state"`
	PostalCode string `bson:"postalCode"`
}

type Contact struct {
	Phone string `bson:"phone"`
	Email string `bson:"email"`
}

type Meter struct {
	MeterID  string      `bson:"meterId"`
	Location GeoJSON     `bson:"location"`
	SIM      SIM         `bson:"sim"`
	Usage    []Usage     `bson:"usage"`
	Alerts   Alert       `bson:"alerts"`
	Commands Commands    `bson:"commands"`
	Status   MeterStatus `bson:"status"`
}

type GeoJSON struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type SIM struct {
	ICCID          string  `bson:"iccid"`
	MobileNumber   string  `bson:"mobileNumber"`
	DataUsageMb    float64 `bson:"dataUsage"`
	ActivationDate int64   `bson:"activationDate"`
}

type Usage struct {
	Date           int64   `bson:"date"`
	ConsumptionKwh float64 `bson:"consumption"`
}

type Alert struct {
	Current CurrentAlert `bson:"current"`
	History []AlertEvent `bson:"history"`
}

type CurrentAlert struct {
	Outage AlertStatus `bson:"outage"`
	Tamper AlertStatus `bson:"tamper"`
}

type AlertStatus struct {
	Active bool   `bson:"active"`
	Since  *int64 `bson:"since,omitempty"`
}

type AlertEvent struct {
	Type      string `bson:"type"`
	StartDate int64  `bson:"startDate"`
	EndDate   int64  `bson:"endDate"`
	Resolved  bool   `bson:"resolved"`
}

type Commands struct {
	Active  []ActiveCommand  `bson:"active"`
	History []HistoryCommand `bson:"history"`
}

type ActiveCommand struct {
	CommandID  string                 `bson:"commandId"`
	Type       CommandType            `bson:"type"`
	IssuedAt   int64                  `bson:"issuedAt"`
	Parameters map[string]interface{} `bson:"parameters"`
	Status     CommandStatus          `bson:"status"`
}

type HistoryCommand struct {
	ActiveCommand `bson:",inline"`
	CompletedAt   int64  `bson:"completedAt"`
	Response      string `bson:"response"`
}

type MeterStatus struct {
	LastSeen       int64   `bson:"lastSeen"`
	GridConnection bool    `bson:"gridConnection"`
	BatteryLevel   float64 `bson:"batteryLevel"`
}

type Billing struct {
	CurrentBill    CurrentBill `bson:"currentBill"`
	PaymentHistory []Payment   `bson:"paymentHistory"`
}

type CurrentBill struct {
	BillingPeriod    DateRange       `bson:"billingPeriod"`
	DueDate          int64           `bson:"dueDate"`
	TotalConsumption decimal.Decimal `bson:"totalConsumption"`
	AmountDue        decimal.Decimal `bson:"amountDue"`
	Paid             bool            `bson:"paid"`
}

type DateRange struct {
	Start int64 `bson:"start"`
	End   int64 `bson:"end"`
}

type Payment struct {
	BillingPeriod DateRange       `bson:"billingPeriod"`
	AmountPaid    decimal.Decimal `bson:"amountPaid"`
	PaymentDate   int64           `bson:"paymentDate"`
	PaymentMethod PaymentMethod   `bson:"paymentMethod"`
	TransactionID string          `bson:"transactionId"`
}

type Ledger struct {
	CurrentBalance Balance      `bson:"currentBalance"`
	UpcomingBill   UpcomingBill `bson:"upcomingBill"`
}

type Balance struct {
	Amount    decimal.Decimal `bson:"amount"`
	UpdatedAt int64           `bson:"updatedAt"`
}

type UpcomingBill struct {
	EstimatedAmount decimal.Decimal `bson:"estimatedAmount"`
	ProjectionDate  int64           `bson:"projectionDate"`
}
