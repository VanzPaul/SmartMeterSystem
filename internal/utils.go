package internal

import (
	"SmartMeterSystem/internal/models"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetResolvedIP() string {
	// Print the address where the server is running
	// Resolve the hostname and IP address
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	var resolvedIP string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			resolvedIP = ipv4.String()
			break
		}
	}
	return resolvedIP
}

// InitLogger initializes a zap logger based on the environment (e.g., production or development)
func NewLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	// Create logger
	isDevelopment, parseErr := strconv.ParseBool(os.Getenv("DEBUG"))
	if parseErr != nil {
		panic(parseErr)
	}

	if !isDevelopment {
		// Use NewProduction for production environments, which includes sensible defaults like JSON format and log rotation.
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err // Return the error if logger creation fails
	}

	return logger, nil
}

func CalculateCharges(doc models.RatesDocument, usedKwh float64, exemptSections, exemptItems map[string]bool) models.Charges {
	log.Printf("=== Calculating Charges ===")
	log.Printf("Used kWh: %.2f, Exempt Sections: %v, Exempt Items: %v", usedKwh, exemptSections, exemptItems)
	log.Printf("Rates Document: Type=%s, Date=%s, Sections=%d",
		doc.Type, doc.RatesData.Date, len(doc.RatesData.Sections))

	charges := models.Charges{
		UsedKwH:   usedKwh,
		AmountDue: 0,
		Rates: models.PaymentRates{
			Date:     doc.RatesData.Date,
			Sections: []models.PaymentSection{}, // Ensure empty slice, not nil
		},
	}

	if len(doc.RatesData.Sections) == 0 {
		log.Println("WARNING: No sections found in rates document")
	}

	for i, section := range doc.RatesData.Sections {
		sectionExempt := exemptSections[section.Name]
		log.Printf("Section %d: %s (Type=%s, Exempt=%v, Items=%d)",
			i, section.Name, section.Type, sectionExempt, len(section.Items))

		if sectionExempt || section.Items == nil {
			log.Printf("SKIPPING SECTION: Exempt=%v, NilItems=%v", sectionExempt, section.Items == nil)
			continue
		}

		paymentSection := models.PaymentSection{
			Name:  section.Name,
			Total: 0,
			Items: []models.PaymentItem{},
		}

		for j, item := range section.Items {
			itemExempt := exemptItems[item.Name]
			log.Printf("  Item %d: %s (Unit=%s, Rate=%.2f, Exempt=%v)",
				j, item.Name, item.Unit, item.Rate, itemExempt)

			if itemExempt {
				continue
			}

			value := 0.0
			switch item.Unit {
			case "PhP/kWh":
				value = usedKwh * item.Rate
			case "PhP/Cust/Mo":
				value = item.Rate
			case "%":
				value = usedKwh * (item.Rate / 100)
			default:
				log.Printf("UNKNOWN UNIT: %s for item %s", item.Unit, item.Name)
			}

			paymentSection.Items = append(paymentSection.Items, models.PaymentItem{
				Name:  item.Name,
				Unit:  item.Unit,
				Rate:  item.Rate,
				Value: value,
			})
			paymentSection.Total += value
		}

		if len(paymentSection.Items) > 0 {
			charges.Rates.Sections = append(charges.Rates.Sections, paymentSection)
			charges.AmountDue += paymentSection.Total
			log.Printf("ADDED SECTION: %s, Total=%.2f", section.Name, paymentSection.Total)
		}
	}

	log.Printf("FINAL CHARGES: AmountDue=%.2f, Sections=%d\n", charges.AmountDue, len(charges.Rates.Sections))
	return charges
}

func GenerateUUIDBillingID() string {
	id := uuid.New() // requires "github.com/google/uuid"
	return fmt.Sprintf("CAL%s", id.String()[:8])
}
