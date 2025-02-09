package models

import "github.com/vanspaul/SmartMeterSystem/models/client"

type Config struct {
	MongoURI string `env:"MONGO_URI"`
	DBName   string `env:"DB_NAME"`
}

type Collection string

const (
	Meters      Collection = "meters"
	Accountings Collection = "accountings"
	Cosumers    Collection = "consumers"
	Accounts    Collection = "accounts"
)

var collectionStructMap = map[Collection]interface{}{
	Accounts:    &Account{},
	Cosumers:    &Consumer{},
	Meters:      &client.Meter{},
	Accountings: &client.Accounting{},
}

func GetCcollectionMap() map[Collection]interface{} {
	return collectionStructMap
}
