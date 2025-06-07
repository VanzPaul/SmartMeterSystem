// internal/server/background.go
package server

import (
	"SmartMeterSystem/internal"
	"SmartMeterSystem/internal/database"
	"SmartMeterSystem/internal/models"
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// BackgroundManager handles all background tasks and cron jobs
type BackgroundManager struct {
	cron   *cron.Cron
	logger *zap.Logger
	server *Server // Reference to server for accessing services if needed
	once   sync.Once
}

// NewBackgroundManager creates a new background tasks manager
func NewBackgroundManager(s *Server) *BackgroundManager {
	return &BackgroundManager{
		cron:   cron.New(cron.WithSeconds()),
		logger: s.logger.Named("background"),
		server: s,
	}
}

// Start initializes and runs all background tasks
func (bm *BackgroundManager) Start() {
	bm.logger.Info("Starting background tasks manager")

	// Register all job groups
	// bm.registerSessionJobs()
	// bm.registerReportJobs()
	// bm.registerMaintenanceJobs()
	bm.issueBillNotice()

	// Start cron scheduler
	bm.cron.Start()

	bm.logger.Info("Background tasks started",
		zap.Int("job_count", len(bm.cron.Entries())))
}

// Stop gracefully shuts down background tasks
func (bm *BackgroundManager) Stop() {
	bm.logger.Info("Stopping background tasks...")
	ctx := bm.cron.Stop()

	// Wait for running jobs to finish
	select {
	case <-ctx.Done():
		bm.logger.Info("Background tasks stopped gracefully")
	case <-time.After(30 * time.Second):
		bm.logger.Warn("Background tasks shutdown timed out")
	}
}

// // Job Group: Session-related tasks
// func (bm *BackgroundManager) registerSessionJobs() {
// 	_, err := bm.cron.AddFunc("@every 10m", bm.cleanupExpiredSessions)
// 	if err != nil {
// 		bm.logger.Error("Failed to register session cleanup job", zap.Error(err))
// 	}

// 	_, err = bm.cron.AddFunc("0 3 * * *", bm.rotateSessionKeys) // Daily at 3AM
// 	if err != nil {
// 		bm.logger.Error("Failed to register session key rotation job", zap.Error(err))
// 	}
// }

// func (bm *BackgroundManager) cleanupExpiredSessions() {
// 	bm.logger.Info("Cleaning up expired sessions")
// 	start := time.Now()

// 	// Example implementation
// 	// err := bm.server.sessionService.CleanupExpired(context.Background())
// 	// if err != nil {
// 	//     bm.logger.Error("Session cleanup failed", zap.Error(err))
// 	// }

// 	bm.logger.Info("Session cleanup completed",
// 		zap.Duration("duration", time.Since(start)))
// }

// func (bm *BackgroundManager) rotateSessionKeys() {
// 	bm.logger.Info("Rotating session keys")
// 	// Implementation here
// }

// Job Group: Report generation
// func (bm *BackgroundManager) registerReportJobs() {
// 	_, err := bm.cron.AddFunc("0 2 * * *", bm.generateNightlyReports) // Daily at 2AM
// 	if err != nil {
// 		bm.logger.Error("Failed to register nightly reports job", zap.Error(err))
// 	}

// 	_, err = bm.cron.AddFunc("0 0 * * 1", bm.generateWeeklyReports) // Monday at 00:00
// 	if err != nil {
// 		bm.logger.Error("Failed to register weekly reports job", zap.Error(err))
// 	}
// }

// -------

func (bm *BackgroundManager) issueBillNotice() {
	_, err := bm.cron.AddFunc("*/30 * * * * *", func() {
		bm.once.Do(func() {

			bm.logger.Info("Issuing Bill Notice")

			const ratesType = "RATES"
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			svc := database.New()
			bm.logger.Debug("Fetching rates document")

			// 1) Fetch the rates document
			var ratesData models.RatesDocument
			findOneResult := svc.FindOne(ctx, "rates", bson.M{"type": ratesType})
			if findOneResult.Err() != nil {
				if findOneResult.Err() == mongo.ErrNoDocuments {
					bm.logger.Sugar().Errorf("Rates document not found (type=%q)", ratesType)
					return
				}
				bm.logger.Sugar().Errorf("Error fetching rates document: %v", findOneResult.Err())
				return
			}
			if err := findOneResult.Decode(&ratesData); err != nil {
				bm.logger.Sugar().Errorf("Failed to decode rates document: %v", err)
				return
			}
			bm.logger.Debug("Rates document loaded",
				zap.String("date", ratesData.RatesData.Date),
				zap.Int("sections", len(ratesData.RatesData.Sections)))

			// 2) Find all active balances
			bm.logger.Debug("Fetching active consumer balances")
			cursorFind, errFindMany := svc.FindMany(ctx, "balances", bson.M{"isActive": true}) // Fixed field name
			if errFindMany != nil {
				bm.logger.Sugar().Errorf("Error fetching active balances: %v", errFindMany)
				return
			}
			defer cursorFind.Close(ctx) // Safe even if cursorFind is nil

			exemptSections := map[string]bool{}
			exemptItems := map[string]bool{}

			// Track processed count for logging
			processedCount := 0
			errorCount := 0

			// 3) Process each active consumer
			for cursorFind.Next(ctx) {
				var balance models.ConsumerBalanceDocument
				if errDecode := cursorFind.Decode(&balance); errDecode != nil {
					errorCount++
					bm.logger.Sugar().Errorf("Failed to decode balance: %v", errDecode)
					continue
				}

				processedCount++
				bm.logger.Debug("Processing consumer",
					zap.Int("id", balance.ID),
					zap.Int("account", balance.AccountNumber),
					zap.String("meter", "MTR-00123")) // Fixed field name

				// Calculate charges with actual consumption (placeholder)
				consumption := 64.00 // Should be actual usage data
				charges := internal.CalculateCharges(ratesData, consumption, exemptSections, exemptItems)
				bm.logger.Debug("Charges calculated",
					zap.Float64("amount", charges.AmountDue),
					zap.Int("sections", len(charges.Rates.Sections)))

				// Generate valid bill with proper dates
				now := time.Now()
				currentBill := models.Billing{
					BillId:    internal.GenerateUUIDBillingID(), // Must be implemented
					IssueDate: now,
					DueDate:   now.AddDate(0, 0, 30), // 30 days from now
					Duration: models.UsageDuration{
						Start: now.AddDate(0, 0, -30), // Last 30 days
						End:   now,
					},
					ConsumerType: "RESIDENTIAL",
					MeterNumber:  "MTR-00123",
					Charges:      charges,
					IsPaid:       false,
				}

				// 4) Update consumer's current bill
				updateFilter := bson.M{"_id": balance.ID}
				updateData := bson.M{"$set": bson.M{"currentBill": currentBill}}
				updateResult, errUpdate := svc.UpdateOne(ctx, "balances", updateFilter, updateData)

				if errUpdate != nil {
					errorCount++
					bm.logger.Sugar().Errorf("Update failed for ID=%v: %v", balance.ID, errUpdate)
				} else if updateResult.MatchedCount == 0 {
					errorCount++
					bm.logger.Sugar().Errorf("No document found for ID=%v", balance.ID)
				} else {
					bm.logger.Debug("Bill updated successfully",
						zap.Int("id", balance.ID),
						zap.Int64("matched", updateResult.MatchedCount),
						zap.Int64("modified", updateResult.ModifiedCount))
				}
			}

			// Final processing summary
			bm.logger.Info("Bill processing completed",
				zap.Int("processed", processedCount),
				zap.Int("errors", errorCount))

			if errCursor := cursorFind.Err(); errCursor != nil {
				bm.logger.Sugar().Errorf("Cursor error: %v", errCursor)
			}
		})

	})

	if err != nil {
		bm.logger.Error("Failed to register cron job", zap.Error(err))
	}

}

// func (bm *BackgroundManager) generateNightlyReports() {
// 	bm.logger.Info("Generating nightly reports")
// 	// Implementation here
// }

// func (bm *BackgroundManager) generateWeeklyReports() {
// 	bm.logger.Info("Generating weekly reports")
// 	// Implementation here
// }

// // Job Group: System maintenance
// func (bm *BackgroundManager) registerMaintenanceJobs() {
// 	_, err := bm.cron.AddFunc("0 4 * * *", bm.cleanupTempFiles) // Daily at 4AM
// 	if err != nil {
// 		bm.logger.Error("Failed to register temp file cleanup job", zap.Error(err))
// 	}

// 	_, err = bm.cron.AddFunc("30 3 * * 6", bm.databaseBackup) // Saturday at 3:30 AM
// 	if err != nil {
// 		bm.logger.Error("Failed to register database backup job", zap.Error(err))
// 	}
// }

// func (bm *BackgroundManager) cleanupTempFiles() {
// 	bm.logger.Info("Cleaning up temporary files")
// 	// Implementation here
// }

// func (bm *BackgroundManager) databaseBackup() {
// 	bm.logger.Info("Performing database backup")
// 	// Implementation here
// }

// Accessing server resources:
// func (bm *BackgroundManager) processPayments() {
// Access database through server reference
//     err := bm.server.db.ProcessPayments()
//     if err != nil {
//         bm.logger.Error("Payment processing failed", zap.Error(err))
//     }
// }

// Handling long-running jobs
// func (bm *BackgroundManager) generateReports() {
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
//     defer cancel()

//     err := bm.server.reportService.Generate(ctx)
// ... error handling ...
// }
