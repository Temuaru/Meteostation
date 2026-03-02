package resource

import (
	"time"
)

// Определение структуры для данных с датчика
type CPUData struct {
	Frequencies []string `json:"frequencies"`
	Unit        string   `json:"unit"`
}

type FrequencyDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Frequency float64   `json:"frequency"`
}

type SensorData struct {
	Temp float32 `json:"temp"`
	Hum  float32 `json:"hum"`
	Co2  uint16  `json:"co2"`
	Pres float32 `json:"pres"`
	Tvoc uint16  `json:"tvoc"`
}

type MeteoDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Co2         uint16    `json:"co2"`
	Tvoc        uint16    `json:"tvoc"`
}
