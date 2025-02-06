package client

import (
	"github.com/shopspring/decimal"
)

// Enums for type safety
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// Accounting
type Accounting struct {
	Billing Billing `bson:"billing"`
	Ledger  Ledger  `bson:"ledger"`
}

// Billing related structures
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

// Ledger related Structures
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
