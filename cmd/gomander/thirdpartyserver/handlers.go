package thirdpartyserver

import (
	"encoding/json"
	"net/http"

	"gomander/internal/command/domain"
	domain2 "gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

func (s *ThirdPartyIntegrationsServer) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"app": "Gomander"}`))
	if err != nil {
		println("Error writing response:", err.Error())
	}
}

func (s *ThirdPartyIntegrationsServer) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commands, err := s.useCases.GetCommands.Execute()
	if err != nil {
		http.Error(w, "Failed to get commands", http.StatusInternalServerError)
		return
	}
	runningCommandsIds := s.useCases.GetRunningCommandIds.Execute()

	mappedCommands := array.Map(commands, func(cmd domain.Command) map[string]interface{} {
		status := "stopped"
		if array.Contains(runningCommandsIds, cmd.Id) {
			status = "running"
		}

		return map[string]interface{}{
			"id":     cmd.Id,
			"name":   cmd.Name,
			"status": status,
		}
	})

	err = json.NewEncoder(w).Encode(mappedCommands)
	if err != nil {
		http.Error(w, "Failed to encode commands", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *ThirdPartyIntegrationsServer) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
}

func (s *ThirdPartyIntegrationsServer) handleStopCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
}

func (s *ThirdPartyIntegrationsServer) handleGetCommandGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groups, err := s.useCases.GetCommandGroups.Execute()
	if err != nil {
		http.Error(w, "Failed to get command groups", http.StatusInternalServerError)
		return
	}
	runningGroupsIds := s.useCases.GetRunningCommandIds.Execute()

	mappedGroups := array.Map(groups, func(group domain2.CommandGroup) map[string]interface{} {
		return map[string]interface{}{
			"id":       group.Id,
			"name":     group.Name,
			"commands": len(group.Commands),
			"runningCommands": len(array.Filter(group.Commands, func(cmd domain.Command) bool {
				return array.Contains(runningGroupsIds, cmd.Id)
			})),
		}
	})

	err = json.NewEncoder(w).Encode(mappedGroups)
	if err != nil {
		http.Error(w, "Failed to encode command groups", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *ThirdPartyIntegrationsServer) handleRunCommandGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
}

func (s *ThirdPartyIntegrationsServer) handleStopCommandGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
}
