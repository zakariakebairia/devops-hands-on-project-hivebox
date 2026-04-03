package main

import (
	"fmt"
	"net/http"
)

const Version = "v1.0.0"

func main() {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, Version)
	})
	http.ListenAndServe(":8080", nil)
}
