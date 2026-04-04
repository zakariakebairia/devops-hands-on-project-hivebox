package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION   = "v1.0.0"
	APIURL    = "https://api.opensensemap.org"
	BOXAPIURL = "https://api.opensensemap.org/boxes"
)

// var client = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════════
// Box struct
// ═══════════════════════════════════════════════════════════════════════
type Box struct {
	ID              string            `json:"_id"`
	Name            string            `json:"name"`
	Location        map[string]string `json:"location"`
	Exposure        string            `json:"exposure"`
	Model           string            `json:"model"`
	Sensors         []Sensor          `json:"sensors"`
	CurrentLocation *Location         `json:"currentLocation"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	Timestamp   time.Time `json:"timestamp"`
}
type Sensor struct {
	ID              string       `json:"_id"`
	Title           string       `json:"title"`
	Unit            string       `json:"unit"` // Temperature, Humidity, PM2.5
	SensorType      string       `json:"sensorType"`
	LastMeasurement *Measurement `json:"lastMeasurement"`
	// Measurements    map[string]time.Time `json:"measurements"` // [value]: date
}
type Measurement struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
}

// ═══════════════════════════════════════════════════════════════════════
// GetBox function
// ═══════════════════════════════════════════════════════════════════════

// senseBox IDs
// ------------
// 5eba5fbad46fb8001b799786
// 5c21ff8f919bf8001adf2488
// 5ade1acf223bd80019a1011c
// json format with an ID
// ----------------------
// https://api.opensensemap.org/boxes/57000b8745fd40c8196ad04c?format=json
func GetData(id string) (*Box, error) {
	response, err := http.Get(BOXAPIURL + "/" + id)
	// response, err := client.Get(BOXAPIURL + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("fetching box %q: %w", id, err)
	}
	defer response.Body.Close()

	var box Box
	if err := json.NewDecoder(response.Body).Decode(&box); err != nil {
		return nil, fmt.Errorf("decoding box %q: %w", id, err)
	}
	fmt.Printf("\n📦 Box name: %q, ID (%s)\n", box.Name, box.ID)
	return &box, nil
}

func GetAllData(ids []string) ([]*Box, error) {
	var boxes []*Box
	for _, id := range ids {
		box, err := GetData(id)
		if err != nil {
			log.Fatal(err)
		}
		boxes = append(boxes, box)
	}
	return boxes, nil
}

func average(nums []float64) (float64, error) {
	if len(nums) == 0 {
		return 0, fmt.Errorf("cannot average empty slice")
	}
	var sum float64
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums)), nil
}

// ═══════════════════════════════════════════════════════════════════════
// GetTemperature function
// ═══════════════════════════════════════════════════════════════════════
// Get the current temperature, and ensure that all data is no older than 1 hour
// TODO: data no older than 1 hour
func GetAverageTemperature(box Box) float64 {
	temps := []float64{}
	for _, sensor := range box.Sensors {
		if strings.Contains(sensor.Unit, "°C") {
			value, err := strconv.ParseFloat(sensor.LastMeasurement.Value, 64)
			if err != nil {
				log.Fatal(err)
			}
			temps = append(temps, value)
		}
	}
	temperature, _ := average(temps)

	fmt.Printf("\n%.2f, average: %.2f", temps, temperature)

	return temperature
}
