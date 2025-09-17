package thirdpartyserver

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"gomander/internal/app"
)

var StartPort = 9001
var EndPort = 9100

type ThirdPartyIntegrationsServer struct {
	useCases app.UseCases
	mu       sync.RWMutex
	Server   *http.Server
}

func NewThirdPartyIntegrationsServer(useCases app.UseCases) *ThirdPartyIntegrationsServer {
	return &ThirdPartyIntegrationsServer{
		useCases: useCases,
	}
}

func (s *ThirdPartyIntegrationsServer) RegisterHandlers() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	port, err := findAvailablePort()
	if err != nil {
		return err
	}

	if s.Server != nil {
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

	s.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return nil
}

func (s *ThirdPartyIntegrationsServer) Start() {
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Handle error (log it, etc.)
			println(err.Error())
		}
	}()
}

func (s *ThirdPartyIntegrationsServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Server == nil {
		return nil // Server is not running
	}

	err := s.Server.Close()
	s.Server = nil
	return err
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
