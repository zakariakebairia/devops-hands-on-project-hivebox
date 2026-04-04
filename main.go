package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	result, err := GetData("5c21ff8f919bf8001adf2488")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(result)
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, VERSION)
	})

	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, GetAverageTemperature(*result))
	})
	http.ListenAndServe(":8080", nil)
}
