package client

// Enums for type safety
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

// Meter-related structures
type Meter struct {
	MeterID     string      `bson:"meterId"`
	MeterConfig MeterConfig `bson:"meterConfig"`
	Location    GeoJSON     `bson:"location"`
	SIM         SIM         `bson:"sim"`
	Usage       []Usage     `bson:"usage"`
	Alerts      Alert       `bson:"alerts"`
	Commands    Commands    `bson:"commands"`
	Status      MeterStatus `bson:"status"`
}

type MeterConfig struct {
	Manufacturer string `bson:"manufacturer"`
	Model        string `bson:"model"`
	Phase        int    `bson:"phase"`
	SerialNumber string `bson:"serialNumber"`
	APIKey       string `bson:"apiKey"`
	WifiSSID     string `bson:"wifiSSID"`
	WifiPassword string `bson:"wifiPassword"`
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
	Date int64   `bson:"date"`
	Kwh  float64 `bson:"consumption"`
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
