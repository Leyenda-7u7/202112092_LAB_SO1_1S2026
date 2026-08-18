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

// Respuesta utilizada por el endpoint /health.
type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"VM"`
	Carnet    string `json:"carnet"`
}

// Respuesta utilizada cuando API1 consulta otra API.
type CallResponse struct {
	APIName    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

// Solo necesitamos leer el campo status de la API consultada.
type RemoteHealth struct {
	Status string `json:"status"`
}

// Cliente HTTP con timeout para evitar que API1 quede esperando
// indefinidamente cuando otra API no esté disponible.
var httpClient = &http.Client{
	Timeout: 3 * time.Second,
}

// Envía respuestas en formato JSON.
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error generando JSON: %v", err)
	}
}

// Verifica el /health de otra API.
func checkAPI(baseURL, apiName, targetVM, carnet string) CallResponse {
	errorResponse := CallResponse{
		APIName:    apiName,
		Message:    fmt.Sprintf("ERROR: The %s located on the %s is not working", apiName, targetVM),
		Connection: false,
		Carnet:     carnet,
	}

	// Si todavía no existe una URL configurada, la API se considera no disponible.
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

	// API1 y API2 viven en VM1.
	vmName := "VM1"

	// API2 estará en otro contenedor dentro de esta misma VM.
	api2URL := os.Getenv("API2_URL")

	// API3 todavía no existe. Configuraremos esta URL cuando creemos VM2.
	api3URL := os.Getenv("API3_URL")

	mux := http.NewServeMux()

	// Endpoint principal de salud.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := HealthResponse{
			Status:    "UP",
			Message:   "API1 is Ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			VM:        vmName,
			Carnet:    carnet,
		}

		writeJSON(w, http.StatusOK, response)
	})

	// API1 -> API2.
	pathAPI2 := fmt.Sprintf("/api1/%s/call-api2", carnet)

	mux.HandleFunc(pathAPI2, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := checkAPI(api2URL, "API2", "VM1", carnet)
		writeJSON(w, http.StatusOK, response)
	})

	// API1 -> API3.
	pathAPI3 := fmt.Sprintf("/api1/%s/call-api3", carnet)

	mux.HandleFunc(pathAPI3, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		response := checkAPI(api3URL, "API3", "VM2", carnet)
		writeJSON(w, http.StatusOK, response)
	})

	log.Printf("API1 lista en el puerto 8080")
	log.Printf("Health: /health")
	log.Printf("API2: %s", pathAPI2)
	log.Printf("API3: %s", pathAPI3)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}