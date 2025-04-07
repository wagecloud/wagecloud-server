package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

// CreateDomainRequest represents the request body for creating a domain
type CreateDomainRequest struct {
	Name   string       `json:"name"`
	Memory model.Memory `json:"memory"`
	Cpu    model.Cpu    `json:"cpu"`
	OS     model.OS     `json:"os"`
}

// CreateDomain handles the creation of a new domain
func (h *Handler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var req CreateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	domain := model.NewDomain(
		model.WithDomainName(req.Name),
		model.WithDomainMemory(req.Memory.Value, req.Memory.Unit),
		model.WithDomainCpu(req.Cpu.Value),
		model.WithDomainOS(model.OS{
			Arch: req.OS.Arch,
			Name: req.OS.Name,
		}),
	)

	_, err := h.service.Libvirt.CreateDomain(domain)
	if err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create domain: "+err.Error())
		return
	}

	response.FromMessage(w, http.StatusCreated, "Domain created successfully")
}

// StartDomain handles starting a domain
func (h *Handler) StartDomain(w http.ResponseWriter, r *http.Request) {
	domainID := chi.URLParam(r, "domainID")
	if domainID == "" {
		response.FromError(w, http.StatusBadRequest, "Domain ID is required")
		return
	}

	// Here you would typically retrieve the domain from libvirt using the ID
	// For now, we'll just return a mock response
	response.FromMessage(w, http.StatusOK, "Domain started successfully")
}
