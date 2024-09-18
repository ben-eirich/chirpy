package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello world")
	mux := http.NewServeMux()
	apiCfg := apiConfig{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerAdminMetrics)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidateChirp)

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

func handlerHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		fmt.Printf("After: %d\n", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerAdminMetrics(w http.ResponseWriter, _ *http.Request) {
	template := `
	<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
	`
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(template, cfg.fileserverHits)))
}

func isBadWord(word string) bool {
	lowerWord := strings.ToLower(word)
	if lowerWord == "kerfuffle" {
		return true
	} else if lowerWord == "sharbert" {
		return true
	} else if lowerWord == "fornax" {
		return true
	} else {
		return false
	}
}

func cleanChirp(input string) string {
	words := strings.Split(input, " ")
	result := []string{}
	for _, word := range words {
		if isBadWord(word) {
			result = append(result, "****")
		} else {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func (cfg *apiConfig) handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Body string `json:"body"`
	}
	type retValue struct {
		Body string `json:"cleaned_body"`
	}

	apiParams := params{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&apiParams)
	if err != nil {
		respondWithError(w, 500, "Error decoding request")
	} else if len(apiParams.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJSON(w, 200, retValue{Body: cleanChirp(apiParams.Body)})
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type retError struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	eo, _ := json.Marshal(retError{Error: msg})
	w.Write(eo)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
	} else {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
