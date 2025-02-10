package models

import (
	"github.com/vanspaul/SmartMeterSystem/models/client"
)

type Collection string

const (
	Meters      Collection = "meters"
	Accountings Collection = "accountings"
	Consumers   Collection = "consumers"
	Accounts    Collection = "accounts"
)

var collectionStructMap = map[Collection]interface{}{
	Accounts:    &Account{},
	Consumers:   &client.Consumer{},
	Meters:      &client.Meter{},
	Accountings: &client.Accounting{},
}

func GetCcollectionMap() map[Collection]interface{} {
	return collectionStructMap
}
