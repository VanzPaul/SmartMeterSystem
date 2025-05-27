// models/meter.go
package models

import (
	"time"
)

type MeterDocument struct {
	ID                    int       `bson:"_id"`
	InstallationDate      time.Time `bson:"meter-installation-date"`
	ConsumerTransformerID string    `bson:"consumer-transformer-id"`
	Latitude              float64   `bson:"meter-latitude"`
	Longitude             float64   `bson:"meter-longitude"`
	ConsumerAccNo         int       `bson:"consumer-acc-no"`
}
