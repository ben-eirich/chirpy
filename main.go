package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello world")
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
