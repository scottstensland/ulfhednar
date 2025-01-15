package main

import (
	"log"
	"net/http"
)

func main() {
	// Define the directory you want to serve
	directory := "./" // Replace with the directory you want to serve

	// Create a file server
	fileServer := http.FileServer(http.Dir(directory))

	// Handle requests and serve the directory
	http.Handle("/", http.StripPrefix("/", fileServer))

	// Start the HTTP server
	port := ":8080" // You can use any port you prefer
	log.Printf("Serving %s on HTTP port %s\n", directory, port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
