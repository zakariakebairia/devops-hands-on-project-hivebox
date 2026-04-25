package main

import (
	"fmt"
	"log"
	"net/http"
)

var Version = "None"

func main() {
	ids := []string{
		"5eba5fbad46fb8001b799786",
		"5c21ff8f919bf8001adf2488",
		"5ade1acf223bd80019a1011c",
	}
	// version endpoint
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, Version)
	})

	// temperature endpoint
	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		boxes, err := FetchBoxesData(ids)
		if err != nil {
			log.Fatal(err)
		}
		temp, err := GetAverageTemperature(boxes)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		}

		fmt.Fprintf(w, `{"average temperature": %.2f}`, temp)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
