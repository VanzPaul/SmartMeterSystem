package main

import (
	"context"
	"log"
	"time"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/controllers"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load environment variables and initialize the logger
	if err := config.LoadEnv(); err != nil {
		// Use standard log if LoadEnv fails before initializing the logger
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	defer utils.Logger.Sync()

	utils.Logger.Debug("debug log") // This will now work if DEBUG=true

	// Rest of your code
	db, err := controllers.NewMongoDB(ctx, &config.MongoEnv)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// var account models.Account
	// filter := bson.M{"email": "johndoe@mail.com"}
	// acc := db.FindOne(ctx, models.Accounts, filter, &account)
	// log.Println(acc)
	// log.Println(account.ID)
	// log.Println(account.HashedPassword)
	// log.Println(account.Email)
	// log.Println(account.Role)

	var accounts []models.Account
	filter := bson.M{"email": "testuser@mail.com"}
	errs := db.Find(ctx, models.Accounts, filter, &accounts)
	log.Println("errs: ", errs)
	log.Println("accounts: ", accounts)
	if len(accounts) == 0 {
		log.Println("zero")
	} else {
		log.Println(len(accounts))
	}

}
