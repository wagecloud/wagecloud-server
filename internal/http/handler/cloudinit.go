package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/internal/http/response"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

// CreateCloudinitRequest represents the request body for creating a cloudinit ISO
type CreateCloudinitRequest struct {
	Userdata      model.Userdata      `json:"userdata"`
	Metadata      model.Metadata      `json:"metadata"`
	NetworkConfig model.NetworkConfig `json:"network_config"`
}

// CreateCloudinit handles the creation of a new cloudinit ISO
func (h *Handler) CreateCloudinit(w http.ResponseWriter, r *http.Request) {
	var req CreateCloudinitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userdata, _ := os.Open("cloud-init-files/ubuntu/user-data")
	metadata, _ := os.Open("cloud-init-files/ubuntu/meta-data")
	networkConfig, _ := os.Open("cloud-init-files/ubuntu/network-config")

	// if err := h.service.Cloudinit.CreateCloudinitByReader(req.Userdata, req.Metadata, req.NetworkConfig); err != nil {
	// 	response.FromError(w, http.StatusInternalServerError, "Failed to create cloudinit ISO: "+err.Error())
	// 	return
	// }

	filename := fmt.Sprintf("cloudinit_%s.iso", uuid.New().String())

	if err := h.service.Cloudinit.CreateCloudinitByReader(filename, userdata, metadata, networkConfig); err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create cloudinit ISO: "+err.Error())
		return
	}

	response.FromMessage(w, http.StatusCreated, "Cloudinit ISO created successfully")
}
