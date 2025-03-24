package handler

import (
	"encoding/json"
	"net/http"

	"github.com/wagecloud/wagecloud-server/internal/http/response"
)

// CreateImageRequest represents the request body for creating an image
type CreateImageRequest struct {
	BaseImagePath  string `json:"baseImagePath"`
	CloneImagePath string `json:"cloneImagePath"`
	Size           uint   `json:"size"`
}

// CreateImage handles the creation of a new image
func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.service.Qemu.CreateImage(req.BaseImagePath, req.CloneImagePath, req.Size)
	if err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create image: "+err.Error())
		return
	}

	response.FromMessage(w, http.StatusCreated, "Image created successfully")
}
