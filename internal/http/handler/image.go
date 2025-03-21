package handler

import (
	"encoding/json"
	"net/http"
)

// CreateImageRequest represents the request body for creating an image
type CreateImageRequest struct {
	BaseImagePath  string `json:"baseImagePath"`
	CloneImagePath string `json:"cloneImagePath"`
}

// CreateImageResponse represents the response for creating an image
type CreateImageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateImage handles the creation of a new image
func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.service.Qemu.ImageCreate(req.BaseImagePath, req.CloneImagePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create image: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateImageResponse{
		Success: true,
		Message: "Image created successfully",
	})
}
