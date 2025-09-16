package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"gomander/internal/app"
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

var StartPort = 9001
var EndPort = 9100

type ThirdPartyIntegrationsServer struct {
	useCases app.UseCases
	mu       sync.RWMutex
	server   *http.Server
}

func NewThirdPartyIntegrationsServer(useCases app.UseCases) *ThirdPartyIntegrationsServer {
	return &ThirdPartyIntegrationsServer{
		useCases: useCases,
	}
}

func (s *ThirdPartyIntegrationsServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	port, err := findAvailablePort()
	if err != nil {
		return err
	}

	if s.server != nil {
		return nil // Server is already running
	}

	mux := http.NewServeMux()

	// Discovery endpoint
	mux.HandleFunc("/discovery", s.handleDiscovery)

	// Commands and Command Groups endpoints
	mux.HandleFunc("/commands", s.handleGetCommands)
	mux.HandleFunc("/command/run/{id}", s.handleRunCommand)
	mux.HandleFunc("/command/stop/{id}", s.handleStopCommand)
	//
	mux.HandleFunc("/command-groups", s.handleGetCommandGroups)
	mux.HandleFunc("/command-group/run/{id}", s.handleRunCommandGroup)
	mux.HandleFunc("/command-group/stop/{id}", s.handleStopCommandGroup)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Handle error (log it, etc.)
			println(err.Error())
		}
	}()

	return nil
}

func (s *ThirdPartyIntegrationsServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.server == nil {
		return nil // Server is not running
	}

	err := s.server.Close()
	s.server = nil
	return err
}

func (s *ThirdPartyIntegrationsServer) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	// Example response for discovery endpoint
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"app": "Gomander"}`))
	if err != nil {
		println("Error writing response:", err.Error())
	}
}

func (s *ThirdPartyIntegrationsServer) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	commands, err := s.useCases.GetCommands.Execute()
	if err != nil {
		http.Error(w, "Failed to get commands", http.StatusInternalServerError)
		return
	}

	mappedCommands := array.Map(commands, func(cmd commanddomain.Command) map[string]interface{} {
		return map[string]interface{}{
			"id":      cmd.Id,
			"name":    cmd.Name,
			"command": cmd.Command,
		}
	})

	err = json.NewEncoder(w).Encode(mappedCommands)
	if err != nil {
		http.Error(w, "Failed to encode commands", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *ThirdPartyIntegrationsServer) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	// Extract command ID from URL
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Command ID is required", http.StatusBadRequest)
	}

	err := s.useCases.RunCommand.Execute(id)
	if err != nil {
		http.Error(w, "Failed to run command", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ThirdPartyIntegrationsServer) handleStopCommand(w http.ResponseWriter, r *http.Request) {
	// Extract command ID from URL
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Command ID is required", http.StatusBadRequest)
	}

	err := s.useCases.StopCommand.Execute(id)
	if err != nil {
		http.Error(w, "Failed to stop command", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ThirdPartyIntegrationsServer) handleGetCommandGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := s.useCases.GetCommandGroups.Execute()
	if err != nil {
		http.Error(w, "Failed to get command groups", http.StatusInternalServerError)
		return
	}

	mappedGroups := array.Map(groups, func(group domain.CommandGroup) map[string]interface{} {
		return map[string]interface{}{
			"id":   group.Id,
			"name": group.Name,
		}
	})

	err = json.NewEncoder(w).Encode(mappedGroups)
	if err != nil {
		http.Error(w, "Failed to encode command groups", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *ThirdPartyIntegrationsServer) handleRunCommandGroup(w http.ResponseWriter, r *http.Request) {
	// Extract command group ID from URL
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Command Group ID is required", http.StatusBadRequest)
	}

	err := s.useCases.RunCommandGroup.Execute(id)
	if err != nil {
		http.Error(w, "Failed to run command group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ThirdPartyIntegrationsServer) handleStopCommandGroup(w http.ResponseWriter, r *http.Request) {
	// Extract command group ID from URL
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Command Group ID is required", http.StatusBadRequest)
	}

	err := s.useCases.StopCommandGroup.Execute(id)
	if err != nil {
		http.Error(w, "Failed to stop command group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func findAvailablePort() (int, error) {
	for port := StartPort; port <= EndPort; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue // Port is in use
		}
		_ = ln.Close() // Close the listener immediately
		return port, nil
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", StartPort, EndPort)
}
