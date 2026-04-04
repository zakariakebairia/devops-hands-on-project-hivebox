package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION = "v0.0.5"
	APIURL  = "https://api.opensensemap.org"
)

// var client = &http.Client{Timeout: 15 * time.Second}

// Box struct
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
}
type Measurement struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
}

// FetchBoxData function

// https://api.opensensemap.org/boxes/57000b8745fd40c8196ad04c?format=json
func FetchBoxData(id string) (*Box, error) {
	// response, err := client.Get(BOXAPIURL + "/" + id)
	response, err := http.Get(APIURL + "/boxes/" + id)
	if err != nil {
		return nil, fmt.Errorf("fetching box %q: %w", id, err)
	}
	defer response.Body.Close()

	var box Box
	if err := json.NewDecoder(response.Body).Decode(&box); err != nil {
		return nil, fmt.Errorf("decoding box %q: %w", id, err)
	}
	fmt.Printf("📦 Box name: %q, ID (%s)\n", box.Name, box.ID)
	return &box, nil
}

func FetchBoxesData(ids []string) ([]*Box, error) {
	var boxes []*Box
	for _, id := range ids {
		box, err := FetchBoxData(id)
		if err != nil {
			return nil, fmt.Errorf("getting box %s: %w", id, err)
		}
		boxes = append(boxes, box)
	}
	return boxes, nil
}

func isLastHour(t time.Time) bool {
	return t.After(time.Now().Add(-1 * time.Hour))
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

// GetAverageTemperature function
// // Get the current temperature, and ensure that all data is no older than 1 hour
// IDEA: the right way to do this is
// Get all boxes -> filter them based on temperature -> get the average value for the temperature
// Best way is , get all data for boxes (ids)
// -> filter only data for temperature (C)
// -> Get values not older than 1 hour
// -> Calculate the average
func GetAverageTemperature(boxes []*Box) (float64, error) {
	var temps []float64
	for _, box := range boxes {
		for _, sensor := range box.Sensors {
			if !strings.Contains(sensor.Unit, "°C") {
				continue
			}
			if sensor.LastMeasurement == nil {
				continue
			}
			if !isLastHour(sensor.LastMeasurement.CreatedAt) {
				continue
			}
			value, err := strconv.ParseFloat(sensor.LastMeasurement.Value, 64)
			if err != nil {
				continue
			}
			temps = append(temps, value)
		}
	}
	if len(temps) == 0 {
		return 0, fmt.Errorf("no temperature readings in the last hour")
	}

	// --------remove later------------
	temperature, err := average(temps)
	if err != nil {
		return 0, fmt.Errorf("getting avg temperature: %w", err)
	}

	fmt.Printf("\n🌡️ temps: %v, Average temperature: %.2f °C\n", temps, temperature)
	// return average(temps)
	return temperature, nil
}
