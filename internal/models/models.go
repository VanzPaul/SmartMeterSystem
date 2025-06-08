// models/meter.go
package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type MeterDocument struct {
	ID                    int        `bson:"_id"`
	MeterNumber           int        `bson:"meterNumber"`
	InstallationDate      time.Time  `bson:"installDate"`
	ConsumerTransformerID string     `bson:"transformerId"`
	Coordinates           []float64  `bson:"coordinates"`
	ConsumerAccNo         int        `bson:"acctNo"`
	SmartMeter            SmartMeter `bson:"smartMeteromitempty"`
	IsActive              bool       `bson:"isActive"`
}

type SmartMeter struct {
	IsActive              bool             `bson:"isActive"`
	Alert                 []Alert          `json:"Alert,omitempty"`
	UsageKwh              float64          `json:"UsageKwh,omitempty"`
	ReadingHistory30days  []ReadingHistory `bson:"readingHistory30days,omitempty"`
	ReadingHistory24hours []ReadingHistory `bson:"readingHistory24hours,omitempty"`
}

type Alert struct {
	ID        string      `json:"ID"`
	Type      AlertType   `json:"Type"`
	Status    AlertStatus `json:"Status"`
	Timestamp string      `json:"Timestamp"`
}

type ReadingHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

//--

type ConsumerDocument struct {
	ID               int       `bson:"_id"`
	AccountNumber    int       `bson:"acctNum"`
	FirstName        string    `bson:"firstName"`
	MiddleName       string    `bson:"middleName"`
	LastName         string    `bson:"lastName"`
	Suffix           string    `bson:"suffix"`
	BirthDate        time.Time `bson:"birthDate"`
	Province         string    `bson:"province"`
	PostalCode       int       `bson:"postalCode"`
	CityMunicipality string    `bson:"cityMun"`
	Barangay         string    `bson:"barangay"`
	Street           string    `bson:"street"`
	PhoneNumber      int       `bson:"phoneNum"`
	IsActive         bool      `bson:"isActive"`
}

type EmployeeDocument struct {
	ID int `bson:"_id"`
}

// ---
// Payment Section

type ConsumerBalanceDocument struct {
	ID             int              `bson:"_id"`
	AccountNumber  int              `bson:"acctNum"`
	ConsumerType   string           `bson:"consumerType"`
	IsActive       bool             `bson:"isActive"`
	LasPaymentDate time.Time        `bson:"lastPaymentDate,omitempty"`
	CurrentBill    Billing          `bson:"currentBill,omitempty"`
	OverdueBill    []Billing        `bson:"overdueBill,omitempty"`
	BillHistory    []Billing        `bson:"billHistory,omitempty"`
	PaymentHistory []PaymentHistory `bson:"paymentHistory,omitempty"`
}

type Billing struct {
	BillId    string        `bson:"billId"`
	IssueDate time.Time     `bson:"issueDate"`
	DueDate   time.Time     `bson:"dueDate"`
	Duration  UsageDuration `bson:"duration"`
	Charges   Charges       `bson:"charges"`
	IsPaid    bool          `bson:"isPaid"`
}

type UsageDuration struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}

type PaymentHistory struct {
	BillIds []string  `bson:"billIds"`
	Date    time.Time `bson:"date"`
	Amount  float64   `bson:"amount"`
	Status  string    `bson:"status"`
}

type Charges struct {
	AmountDue   float64      `bson:"amountDue"`
	UsedKwH     float64      `bson:"usedKwH"`
	Rates       PaymentRates `bson:"rates"`
	OverdueFees *OverdueFees `bson:"overdueFees,omitempty"`
}

type PaymentRates struct {
	Date     string           `json:"date"`
	Sections []PaymentSection `json:"sections"`
}

type PaymentSection struct {
	Name  string        `json:"name"`
	Total float64       `json:"total"`
	Items []PaymentItem `json:"items,omitempty"`
}

type PaymentItem struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Rate  float64 `json:"rate"`
	Value float64 `json:"value"`
}

type OverdueFees struct {
	ServiceFee OverdueServiceFee `json:"serviceFee"`
	Interest   OverdueInterest   `json:"interest"`
}

type OverdueServiceFee struct {
	Rate  float64 `json:"rate"`
	Value float64 `json:"value"`
}

type OverdueInterest struct {
	Rate  float64 `json:"rate"`
	Value float64 `json:"value"`
}

// ---

type AlertType string

// All of the AlertType values your system supports.
const (
	AlertTypePowerOutage AlertType = "power_outage"
	AlertTypeTamper      AlertType = "tamper"
	AlertTypeLowBattery  AlertType = "low_battery"
)

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusInactive AlertStatus = "inactive"
)

type SmartMeterDocument struct {
	ID                    string           `json:"ID"`
	Number                string           `json:"Number"`
	Location              string           `json:"Location"`
	Latitude              float64          `json:"Latitude"`
	Longitude             float64          `json:"Longitude"`
	Status                string           `json:"Status,omitempty"`
	Alert                 []Alert          `json:"Alert,omitempty"`
	UsageKwh              float64          `json:"UsageKwh,omitempty"`
	ReadingHistory30days  []ReadingHistory `bson:"readingHistory30days,omitempty"`
	ReadingHistory24hours []ReadingHistory `bson:"readingHistory24hours,omitempty"`
}

// type Alert struct {
// 	ID        string      `json:"ID"`
// 	Type      AlertType   `json:"Type"`
// 	Status    AlertStatus `json:"Status"`
// 	Timestamp string      `json:"Timestamp"`
// }

// type ReadingHistory struct {
// 	Timestamp time.Time `json:"timestamp"`
// 	Value     float64   `json:"value"`
// }

// ---
var AccountingRatesTableFormType = struct {
	Display   string
	FormRates string
	FormERC   string
}{
	Display:   "display",
	FormRates: "form-rates",
	FormERC:   "form-erc",
}

type AccountingRatesTable struct {
	Date        string `json:"date,omitempty"`
	Particulars string `json:"particulars"`
	// Unit                          string `json:"unit"`
	Rates                        string                         `json:"rates,,omitempty"`
	ERC                          string                         `json:"erc,,omitempty"`
	AccountingRatesTableRowGroup []AccountingRatesTableRowGroup `json:"row-group"`
}

type AccountingRatesTableRowGroup struct {
	Particulars string        `json:"particulars"`
	Unit        string        `json:"unit,omitempty"`
	Rates       string        `json:"rates,omitempty"`
	ERC         string        `json:"erc,omitempty"`
	SubRowGroup []SubRowGroup `json:"sub-row-group"`
}

type SubRowGroup struct {
	Particulars string `json:"particulars"`
	Unit        string `json:"unit"`
	Rates       string `json:"rates,omitempty"`
	ERC         string `json:"erc,omitempty"`
}

func (art *AccountingRatesTable) UnmarshalJSON(data []byte) error {
	type Alias AccountingRatesTable
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(art),
	}

	// First unmarshal normally to get simple fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Now handle the flattened structure
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Track row groups by index
	rowGroups := make(map[int]AccountingRatesTableRowGroup)

	for key, value := range raw {
		if strings.HasPrefix(key, "row-group[") {
			parts := strings.SplitN(key, "].", 2)
			if len(parts) != 2 {
				continue
			}

			// Extract index
			indexStr := strings.TrimPrefix(parts[0], "row-group[")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}

			// Get or create row group
			rg, exists := rowGroups[index]
			if !exists {
				rg = AccountingRatesTableRowGroup{}
			}

			// Handle sub fields
			fieldPath := parts[1]
			switch {
			case strings.HasPrefix(fieldPath, "sub-row-group["):
				subParts := strings.SplitN(fieldPath, "].", 2)
				if len(subParts) != 2 {
					continue
				}

				// Extract sub index
				subIndexStr := strings.TrimPrefix(subParts[0], "sub-row-group[")
				subIndex, err := strconv.Atoi(subIndexStr)
				if err != nil {
					return err
				}

				// Ensure subrow group exists
				for len(rg.SubRowGroup) <= subIndex {
					rg.SubRowGroup = append(rg.SubRowGroup, SubRowGroup{})
				}

				// Set subrow group field
				field := subParts[1]
				switch field {
				case "particulars":
					rg.SubRowGroup[subIndex].Particulars = value.(string)
				case "unit":
					rg.SubRowGroup[subIndex].Unit = value.(string)
				case "rates":
					rg.SubRowGroup[subIndex].Rates = value.(string)
				case "erc":
					rg.SubRowGroup[subIndex].ERC = value.(string)
				}
			default:
				// Handle top-level row group fields
				switch fieldPath {
				case "particulars":
					rg.Particulars = value.(string)
				case "unit":
					rg.Unit = value.(string)
				case "rates":
					rg.Rates = value.(string)
				case "erc":
					rg.ERC = value.(string)
				}
			}

			rowGroups[index] = rg
		}
	}

	// Convert map to ordered slice
	maxIndex := 0
	for index := range rowGroups {
		if index > maxIndex {
			maxIndex = index
		}
	}

	art.AccountingRatesTableRowGroup = make([]AccountingRatesTableRowGroup, maxIndex+1)
	for index, rg := range rowGroups {
		art.AccountingRatesTableRowGroup[index] = rg
	}

	return nil
}

// ---

// Define data structures matching client-side format

// Rates and ERC

type RatesDocument struct {
	Type         string       `json:"type" bson:"type"`
	RatesData    RatesData    `json:"ratesdata" bson:"ratesdata"`
	InterestData InterestData `json:"interestdata" bson:"interestdata"`
}

type RatesData struct {
	Date     string    `json:"date" bson:"date"`
	Type     string    `json:"type" bson:"type"`
	Sections []Section `json:"sections" bson:"sections"`
}

type Section struct {
	Type  string  `json:"type" bson:"type"`
	Name  string  `json:"name" bson:"name"`
	Rate  float64 `json:"rate" bson:"rate,omitempty"`
	ERC   float64 `json:"erc" bson:"erc,omitempty"`
	Items []Item  `json:"items" bson:"items,omitempty"`
}

type Item struct {
	Name string  `json:"name" bson:"name"`
	Unit string  `json:"unit" bson:"unit"`
	Rate float64 `json:"rate" bson:"rate,omitempty"`
	ERC  float64 `json:"erc" bson:"erc,omitempty"`
}

type InterestData struct {
	Interest float64 `json:"interest" bson:"interest"`
}

var RatesSampleData RatesData = RatesData{
	Date: "01/01/01",
	Type: "RESIDENTIAL",
	Sections: []Section{
		{Type: "main-header", Name: "RESIDENTIAL"},
		{
			Type: "category", Name: "Generation Charges", Rate: 5.6092, ERC: 5.6092,
			Items: []Item{
				{"Generation Energy Charge", "PhP/kWh", 5.6092, 5.6092},
				{"Other Generation Rate Adjustment", "PhP/kWh", 0.0, 0.0},
			},
		},
		{
			Type: "category", Name: "Transmission Charges (NGCP)", Rate: 0.6853, ERC: 0.6853,
			Items: []Item{
				{"Transmission Demand Charge", "PhP/kW", 0.0, 0.0},
				{"Transmission System Charge", "PhP/kWh", 0.6853, 0.6853},
			},
		},
		{
			Type: "category", Name: "System Loss Charge", Rate: 0.9344, ERC: 0.9344,
			Items: []Item{
				{"System Loss Charge", "PhP/kWh", 0.9344, 0.9344},
			},
		},
		{
			Type: "category", Name: "Distribution Charges", Rate: 0.4613, ERC: 0.4613,
			Items: []Item{
				{"Distribution Demand Charge", "PhP/kW", 0.0, 0.0},
				{"Distribution System Charge", "PhP/kWh", 0.4613, 0.4613},
			},
		},
		{
			Type: "category", Name: "Supply Charges", Rate: 0.5376, ERC: 0.5376,
			Items: []Item{
				{"Supply Retail Customer Charge", "PhP/Cust/Mo", 0.0, 0.0},
				{"Supply System Charge", "PhP/kWh", 0.5376, 0.5376},
			},
		},
		{
			Type: "category", Name: "Metering Charges", Rate: 5.3205, ERC: 0.3205,
			Items: []Item{
				{"Metering Retail Customer Charge", "PhP/Cust/Mo", 5.0000, 0.0},
				{"Metering System Charge", "PhP/kWh", 0.3205, 0.3205},
			},
		},
		{
			Type: "category", Name: "Reinvestment Fund/MCC", Rate: 0.2178, ERC: 0.2178,
			Items: []Item{
				{"Reinvestment Fund/MCC", "PhP/kWh", 0.2178, 0.2178},
			},
		},
		{
			Type: "category", Name: "Other Charges", Rate: 0.4506, ERC: 0.0006,
			Items: []Item{
				{"Lifeline Subsidy Charge", "PhP/kWh", 0.0003, 0.0003},
				{"Sr. Citizen Subsidy Charge", "PhP/kWh", 0.0003, 0.0003},
				{"Lifeline Discount (0-15 kWh)", "", 0.4000, 0.0},
				{"Lifeline Discount (16-35 kWh)", "", 0.0500, 0.0},
			},
		},
		{
			Type: "category", Name: "Franchise Tax", Rate: 0.0050, ERC: 0.0,
			Items: []Item{
				{"Franchise Tax", "", 0.0050, 0.0},
			},
		},
		{
			Type: "category", Name: "Business Tax", Rate: 0.0307, ERC: 0.0,
			Items: []Item{
				{"Business Tax - Balayan", "PhP/kWh", 0.0041, 0.0},
				{"Business Tax - Calaca", "PhP/kWh", 0.0117, 0.0},
				{"Business Tax - Calatagan", "PhP/kWh", 0.0040, 0.0},
				{"Business Tax - Lemery", "PhP/kWh", 0.0033, 0.0},
				{"Business Tax - Nasugbu", "PhP/kWh", 0.0034, 0.0},
				{"Business Tax - Tuy", "PhP/kWh", 0.0042, 0.0},
			},
		},
		{
			Type: "category", Name: "Real Property Tax", Rate: 0.0367, ERC: 0.0,
			Items: []Item{
				{"Real Property Tax - Calaca", "PhP/kWh", 0.0219, 0.0},
				{"Real Property Tax - Balayan", "PhP/kWh", 0.0027, 0.0},
				{"Real Property Tax - Lemery", "PhP/kWh", 0.0028, 0.0},
				{"Real Property Tax - Nasugbu", "PhP/kWh", 0.0034, 0.0},
				{"Real Property Tax - Taal", "PhP/kWh", 0.0030, 0.0},
				{"Real Property Tax - Calatagan", "PhP/kWh", 0.0029, 0.0},
			},
		},
		{
			Type: "category", Name: "VAT", Rate: 1.0943, ERC: 0.8543,
			Items: []Item{
				{"Generation", "PhP/kWh", 0.6376, 0.6376},
				{"Transmission", "PhP/kWh", 0.1096, 0.1096},
				{"System Loss", "PhP/kWh", 0.1071, 0.1071},
				{"GRAM/ICERA/DAA VAT", "PhP/kWh", 0.0, 0.0},
				{"Distribution %", "", 0.1200, 0.0},
				{"Others %", "", 0.1200, 0.0},
			},
		},
		{
			Type: "category", Name: "Universal Charge", Rate: 0.2250, ERC: 0.2250,
			Items: []Item{
				{"Missionary Electrification", "PhP/kWh", 0.1822, 0.1822},
				{"True-up (CY 2012)", "PhP/kWh", 0.0, 0.0},
				{"True-up (CY 2013)", "PhP/kWh", 0.0, 0.0},
				{"True-up (CY 2014)", "PhP/kWh", 0.0, 0.0},
				{"Environmental Charge", "PhP/kWh", 0.0, 0.0},
				{"NPC Stranded Contract Cost", "PhP/kWh", 0.0, 0.0},
				{"NPC Stranded Debts", "PhP/kWh", 0.0428, 0.0428},
				{"GRAM/ICERA/DAA", "PhP/kWh", 0.0, 0.0},
			},
		},
		{
			Type: "category", Name: "FIT - ALL", Rate: 0.0838, ERC: 0.0838,
			Items: []Item{
				{"FIT-ALL Php/kWh", "PhP/kWh", 0.0838, 0.0838},
			},
		},
		{
			Type: "total", Name: "TOTAL RATE", Rate: 9.0000, ERC: 9.9298,
		},
	},
}

/*
var ratesData string = `
{
        billingDate: "2023-01-01",
        type: "RESIDENTIAL",
        overdueInterest: 7.29,
        sections: [
            {
                id: "header-residential",
                type: "main-header",
                name: "RESIDENTIAL",
                rate: "",
                erc: 9.9298
            },
            {
                id: "cat-gen",
                type: "category",
                name: "Generation Charges",
                rate: 5.6092,
                erc: 5.6092,
                items: [
                    { id: "item-gen1", name: "Generation Energy Charge", unit: "PhP/kWh", rate: 5.6092, erc: 5.6092 },
                    { id: "item-gen2", name: "Other Generation Rate Adjustment", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 }
                ]
            },
            {
                id: "cat-trans",
                type: "category",
                name: "Transmission Charges (NGCP)",
                rate: 0.6853,
                erc: 0.6853,
                items: [
                    { id: "item-trans1", name: "Transmission Demand Charge", unit: "PhP/kW", rate: 0.0000, erc: 0.0000 },
                    { id: "item-trans2", name: "Transmission System Charge", unit: "PhP/kWh", rate: 0.6853, erc: 0.6853 }
                ]
            },
            {
                id: "cat-sysloss",
                type: "category",
                name: "System Loss Charge",
                rate: 0.9344,
                erc: 0.9344,
                items: [
                    { id: "item-sysloss1", name: "System Loss Charge", unit: "PhP/kWh", rate: 0.9344, erc: 0.9344 }
                ]
            },
            {
                id: "cat-dist",
                type: "category",
                name: "Distribution Charges",
                rate: 0.4613,
                erc: 0.4613,
                items: [
                    { id: "item-dist1", name: "Distribution Demand Charge", unit: "PhP/kW", rate: 0.0000, erc: 0.0000 },
                    { id: "item-dist2", name: "Distribution System Charge", unit: "PhP/kWh", rate: 0.4613, erc: 0.4613 }
                ]
            },
            {
                id: "cat-supply",
                type: "category",
                name: "Supply Charges",
                rate: 0.5376,
                erc: 0.5376,
                items: [
                    { id: "item-supply1", name: "Supply Retail Customer Charge", unit: "PhP/Cust/Mo", rate: 0.0000, erc: 0.0000 },
                    { id: "item-supply2", name: "Supply System Charge", unit: "PhP/kWh", rate: 0.5376, erc: 0.5376 }
                ]
            },
            {
                id: "cat-meter",
                type: "category",
                name: "Metering Charges",
                rate: 5.3205,
                erc: 0.3205,
                items: [
                    { id: "item-meter1", name: "Metering Retail Customer Charge", unit: "PhP/Cust/Mo", rate: 5.0000, erc: 0.0000 },
                    { id: "item-meter2", name: "Metering System Charge", unit: "PhP/kWh", rate: 0.3205, erc: 0.3205 }
                ]
            },
            {
                id: "cat-reinvest",
                type: "category",
                name: "Reinvestment Fund/MCC",
                rate: 0.2178,
                erc: 0.2178,
                items: [
                    { id: "item-reinvest1", name: "Reinvestment Fund/MCC", unit: "PhP/kWh", rate: 0.2178, erc: 0.2178 }
                ]
            },
            {
                id: "cat-other",
                type: "category",
                name: "Other Charges",
                rate: 0.4506,
                erc: 0.0006,
                items: [
                    { id: "item-other1", name: "Lifeline Subsidy Charge", unit: "PhP/kWh", rate: 0.0003, erc: 0.0003 },
                    { id: "item-other2", name: "Sr. Citizen Subsidy Charge", unit: "PhP/kWh", rate: 0.0003, erc: 0.0003 },
                    { id: "item-other3", name: "Lifeline Discount (0-15 kWh)", unit: "", rate: 0.4000, erc: 0.0000 },
                    { id: "item-other4", name: "Lifeline Discount (16-35 kWh)", unit: "", rate: 0.0500, erc: 0.0000 }
                ]
            },
            {
                id: "cat-franchise",
                type: "category",
                name: "Franchise Tax",
                rate: 0.0050,
                erc: 0.0000,
                items: [
                    { id: "item-franchise1", name: "Franchise Tax", unit: "", rate: 0.0050, erc: 0.0000 }
                ]
            },
            {
                id: "cat-business",
                type: "category",
                name: "Business Tax",
                rate: 0.0307,
                erc: 0.0000,
                items: [
                    { id: "item-business1", name: "Business Tax - Balayan", unit: "PhP/kWh", rate: 0.0041, erc: 0.0000 },
                    { id: "item-business2", name: "Business Tax - Calaca", unit: "PhP/kWh", rate: 0.0117, erc: 0.0000 },
                    { id: "item-business3", name: "Business Tax - Calatagan", unit: "PhP/kWh", rate: 0.0040, erc: 0.0000 },
                    { id: "item-business4", name: "Business Tax - Lemery", unit: "PhP/kWh", rate: 0.0033, erc: 0.0000 },
                    { id: "item-business5", name: "Business Tax - Nasugbu", unit: "PhP/kWh", rate: 0.0034, erc: 0.0000 },
                    { id: "item-business6", name: "Business Tax - Tuy", unit: "PhP/kWh", rate: 0.0042, erc: 0.0000 }
                ]
            },
            {
                id: "cat-property",
                type: "category",
                name: "Real Property Tax",
                rate: 0.0367,
                erc: 0.0000,
                items: [
                    { id: "item-property1", name: "Real Property Tax - Calaca", unit: "PhP/kWh", rate: 0.0219, erc: 0.0000 },
                    { id: "item-property2", name: "Real Property Tax - Balayan", unit: "PhP/kWh", rate: 0.0027, erc: 0.0000 },
                    { id: "item-property3", name: "Real Property Tax - Lemery", unit: "PhP/kWh", rate: 0.0028, erc: 0.0000 },
                    { id: "item-property4", name: "Real Property Tax - Nasugbu", unit: "PhP/kWh", rate: 0.0034, erc: 0.0000 },
                    { id: "item-property5", name: "Real Property Tax - Taal", unit: "PhP/kWh", rate: 0.0030, erc: 0.0000 },
                    { id: "item-property6", name: "Real Property Tax - Calatagan", unit: "PhP/kWh", rate: 0.0029, erc: 0.0000 }
                ]
            },
            {
                id: "cat-vat",
                type: "category",
                name: "VAT",
                rate: 1.0943,
                erc: 0.8543,
                items: [
                    { id: "item-vat1", name: "Generation", unit: "PhP/kWh", rate: 0.6376, erc: 0.6376 },
                    { id: "item-vat2", name: "Transmission", unit: "PhP/kWh", rate: 0.1096, erc: 0.1096 },
                    { id: "item-vat3", name: "System Loss", unit: "PhP/kWh", rate: 0.1071, erc: 0.1071 },
                    { id: "item-vat4", name: "GRAM/ICERA/DAA VAT", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-vat5", name: "Distribution %", unit: "", rate: 0.1200, erc: 0.0000 },
                    { id: "item-vat6", name: "Others %", unit: "", rate: 0.1200, erc: 0.0000 }
                ]
            },
            {
                id: "cat-universal",
                type: "category",
                name: "Universal Charge",
                rate: 0.2250,
                erc: 0.2250,
                items: [
                    { id: "item-universal1", name: "Missionary Electrification", unit: "PhP/kWh", rate: 0.1822, erc: 0.1822 },
                    { id: "item-universal2", name: "True-up (CY 2012)", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-universal3", name: "True-up (CY 2013)", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-universal4", name: "True-up (CY 2014)", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-universal5", name: "Environmental Charge", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-universal6", name: "NPC Stranded Contract Cost", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 },
                    { id: "item-universal7", name: "NPC Stranded Debts", unit: "PhP/kWh", rate: 0.0428, erc: 0.0428 },
                    { id: "item-universal8", name: "GRAM/ICERA/DAA", unit: "PhP/kWh", rate: 0.0000, erc: 0.0000 }
                ]
            },
            {
                id: "cat-fit",
                type: "category",
                name: "FIT - ALL",
                rate: 0.0838,
                erc: 0.0838,
                items: [
                    { id: "item-fit1", name: "FIT-ALL Php/kWh", unit: "PhP/kWh", rate: 0.0838, erc: 0.0838 }
                ]
            },
            {
                id: "total-section",
                type: "total",
                name: "TOTAL RATE",
                rate: 9.0000,
                erc: 9.9298
            }
        ]
    };
`
*/

// ---
// Payment

// async function fetchAccountDetails(accountNumber) {
//     return new Promise((resolve, reject) => {
//         setTimeout(() => {
//             if (accountNumber.startsWith('00')) {
//                 reject(new Error('Consumer account number does not exist or is invalid'));
//             } else {
//                 const today = new Date();
//                 const dueDate = new Date();
//                 dueDate.setDate(today.getDate() + 7); // Due in 7 days

//                 const accountData = {
//                     accountNumber: "AC982345",
//                     firstName: "Maria",
//                     middleName: "Santos",
//                     lastName: "Dela Cruz",
//                     suffix: "",
//                     lastPaymentDate: "2023-06-15",
//                     dueDate: dueDate.toISOString().split('T')[0],
// 					currentBill: {
// 						id: <replace with sample id here>,
// 						currentBillMonth: "July",
// 						currentBillYear: "2025",
// 						currentBillAmount: 1850.75,
// 					}
//                     overdueBills: [
//                         {
// 							id: <replace with sample id here>
//                             month: "June",
//                             year: "2025",
//                             baseAmount: 850.25,
//                             interestRate: 3.5,
//                             interestAmount: 29.76,
//                             serviceFee: 50.00,
//                             daysOverdue: 45
//                         },
//                         {
// 							id: <replace with sample id here>
//                             month: "May",
//                             year: "2025",
//                             baseAmount: 750.50,
//                             interestRate: 3.5,
//                             interestAmount: 26.27,
//                             serviceFee: 50.00,
//                             daysOverdue: 75
//                         }
//                     ],
//                     totalAmount: 1850.75 + 850.25 + 29.76 + 50.00 + 750.50 + 26.27 + 50.00,
//                     paymentHistory: [
//                         { date: '2023-06-15', transactionId: 'TX789012', amount: 1750.50, status: 'Paid' },
//                         { date: '2023-05-12', transactionId: 'TX345678', amount: 1650.25, status: 'Paid' },
//                         { date: '2023-04-10', transactionId: 'TX901234', amount: 1800.00, status: 'Paid' },
//                         { date: '2023-03-08', transactionId: 'TX567890', amount: 1720.30, status: 'Paid' }
//                     ]
//                 };
//                 resolve(accountData);
//             }
//         }, 1500);
//     });
// }
