package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
   )
   func handler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	// Get the 'min' and 'max' query parameters from the URL
	minParam := r.URL.Query().Get("min")
	maxParam := r.URL.Query().Get("max")
	// Convert the parameters to integers
	min, _ := strconv.Atoi(minParam)
	max, _ := strconv.Atoi(maxParam)
	// Generate a random number within the specified range
	randomNumber := rand.Intn(max-min+1) + min
	fmt.Fprintf(w, "Random Number: %d", randomNumber)
   }
   func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
   }