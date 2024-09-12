package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello world")
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthzHandler)
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

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
