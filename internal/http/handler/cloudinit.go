package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/cloudinit"
)

// CreateCloudinitRequest represents the request body for creating a cloudinit ISO
type CreateCloudinitRequest struct {
	Userdata model.Userdata `json:"userdata"`
	Metadata model.Metadata `json:"metadata"`
}

// CreateCloudinit handles the creation of a new cloudinit ISO
func (h *Handler) CreateCloudinit(w http.ResponseWriter, r *http.Request) {
	var req CreateCloudinitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	params := cloudinit.CreateCloudinitParams{
		Userdata: req.Userdata,
		Metadata: req.Metadata,
	}

	isoReader, err := h.service.Cloudinit.CreateCloudinit(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create cloudinit ISO: "+err.Error())
		return
	}

	// Set appropriate headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=cloudinit.iso")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	// Stream the ISO file to the client
	_, err = io.Copy(w, isoReader)
	if err != nil {
		// Note: We can't send an error response here as we've already started writing the response
		http.Error(w, "Failed to stream ISO file", http.StatusInternalServerError)
		return
	}
}
