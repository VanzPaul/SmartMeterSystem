// models/meter.go
package models

import (
	"time"
)

type MeterDocument struct {
	ID                    int       `bson:"_id"`
	InstallationDate      time.Time `bson:"installDate"`
	ConsumerTransformerID string    `bson:"transformerId"`
	Latitude              float64   `bson:"lat"`
	Longitude             float64   `bson:"long"`
	ConsumerAccNo         int       `bson:"acctNo"`
}

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
}
