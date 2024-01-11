package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func ExecuteCommand(commandRequest CommandRequest) (string, error) {
	cmd := exec.Command("sh", "-c", commandRequest.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing the command: %v", err)
	}
	return string(output), nil
}

func CommandRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var commandRequest CommandRequest
	fmt.Printf("received command...")
	if err := json.NewDecoder(request.Body).Decode(&commandRequest); err != nil {
		fmt.Printf("Failed to decode request  %v", err)
		http.Error(responseWriter, fmt.Sprintf("Failed to decode request %v", err), http.StatusBadRequest)
		return
	}

	if commandRequest.Command == "" {
		http.Error(responseWriter, "Empty command", http.StatusBadRequest)
		return
	}
	fmt.Printf("Will execute command...")
	output, err := ExecuteCommand(commandRequest)
	fmt.Printf("done executing command...")
	if err != nil {
		response := CommandResponse{
			Error:  err.Error(),
			Output: "",
		}
		json.NewEncoder(responseWriter).Encode(response)
		return
	}
	response := CommandResponse{
		Error:  "",
		Output: output,
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/cmd", CommandRequestHandler).Methods("POST")
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}).Methods("GET")
	http.Handle("/", router)

	fmt.Printf("Will stat server on %d", 8080)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
	if err != nil {
		fmt.Printf("Error statring server: %v", err)
	}
}
