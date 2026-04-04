package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 5ade1acf223bd80019a1011c
	// 5c21ff8f919bf8001adf2488
	result, err := GetData("5c21ff8f919bf8001adf2488")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, VERSION)
	})

	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		temp, err := GetAverageTemperature(*result)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		}

		fmt.Fprintf(w, `{"average temperature": %.2f}`, temp)
	})
	http.ListenAndServe(":8080", nil)
}
