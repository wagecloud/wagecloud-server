package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

// CreateDomainRequest represents the request body for creating a domain
type CreateDomainRequest struct {
	Name          string       `json:"name"`
	UUID          string       `json:"uuid"`
	Memory        model.Memory `json:"memory"`
	Cpu           model.Cpu    `json:"cpu"`
	Arch          string       `json:"arch"`
	SourcePath    string       `json:"sourcePath"`
	CloudinitPath string       `json:"cloudinitPath"`
}

// CreateDomainResponse represents the response for creating a domain
type CreateDomainResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateDomain handles the creation of a new domain
func (h *Handler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var req CreateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	domain := model.Domain{
		Name:          req.Name,
		UUID:          req.UUID,
		Memory:        req.Memory,
		Cpu:           req.Cpu,
		Arch:          model.Arch(req.Arch),
		SourcePath:    req.SourcePath,
		CloudinitPath: req.CloudinitPath,
	}

	_, err := h.service.Libvirt.CreateDomain(domain)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create domain: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateDomainResponse{
		Success: true,
		Message: "Domain created successfully",
	})
}

// StartDomain handles starting a domain
func (h *Handler) StartDomain(w http.ResponseWriter, r *http.Request) {
	domainID := chi.URLParam(r, "domainID")
	if domainID == "" {
		respondWithError(w, http.StatusBadRequest, "Domain ID is required")
		return
	}

	// Here you would typically retrieve the domain from libvirt using the ID
	// For now, we'll just return a mock response
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Domain started successfully",
	})
}
