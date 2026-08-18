package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"VM"`
	Carnet    string `json:"carnet"`
}

type CallResponse struct {
	APIName    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

type RemoteHealth struct {
	Status string `json:"status"`
}

var httpClient = &http.Client{
	Timeout: 3 * time.Second,
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error generando JSON: %v", err)
	}
}

func checkAPI(baseURL, apiName, targetVM, carnet string) CallResponse {
	errorResponse := CallResponse{
		APIName:    apiName,
		Message:    fmt.Sprintf("ERROR: The %s located on the %s is not working", apiName, targetVM),
		Connection: false,
		Carnet:     carnet,
	}

	if strings.TrimSpace(baseURL) == "" {
		return errorResponse
	}

	healthURL := strings.TrimRight(baseURL, "/") + "/health"

	resp, err := httpClient.Get(healthURL)
	if err != nil {
		return errorResponse
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errorResponse
	}

	var health RemoteHealth

	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return errorResponse
	}

	if health.Status != "UP" {
		return errorResponse
	}

	return CallResponse{
		APIName:    apiName,
		Message:    fmt.Sprintf("The %s located on the %s is working", apiName, targetVM),
		Connection: true,
		Carnet:     carnet,
	}
}

func main() {
	carnet := strings.TrimSpace(os.Getenv("CARNET"))

	if carnet == "" {
		log.Fatal("La variable de entorno CARNET es obligatoria")
	}

	vmName := "VM2"

	api1URL := os.Getenv("API1_URL")
	api2URL := os.Getenv("API2_URL")

	mux := http.NewServeMux()

	// GET /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := HealthResponse{
			Status:    "UP",
			Message:   "API3 is Ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			VM:        vmName,
			Carnet:    carnet,
		}

		writeJSON(w, http.StatusOK, response)
	})

	// GET /api3/202112092/call-api1
	pathAPI1 := fmt.Sprintf("/api3/%s/call-api1", carnet)

	mux.HandleFunc(pathAPI1, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := checkAPI(api1URL, "API1", "VM1", carnet)
		writeJSON(w, http.StatusOK, response)
	})

	// GET /api3/202112092/call-api2
	pathAPI2 := fmt.Sprintf("/api3/%s/call-api2", carnet)

	mux.HandleFunc(pathAPI2, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := checkAPI(api2URL, "API2", "VM1", carnet)
		writeJSON(w, http.StatusOK, response)
	})

	log.Printf("API3 lista en el puerto 8080")
	log.Printf("Health: /health")
	log.Printf("API1: %s", pathAPI1)
	log.Printf("API2: %s", pathAPI2)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}