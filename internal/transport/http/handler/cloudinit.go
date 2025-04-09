package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/middleware/auth"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

// CreateCloudinitRequest represents the request body for creating a cloudinit ISO
type CreateCloudinitRequest struct {
	Userdata      model.Userdata      `json:"userdata"`
	Metadata      model.Metadata      `json:"metadata"`
	NetworkConfig model.NetworkConfig `json:"network_config"`
}

// CreateCloudinit handles the creation of a new cloudinit ISO
func (h *Handler) CreateCloudinit(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	var req CreateCloudinitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromHTTPError(w, http.StatusBadRequest)
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
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, nil, http.StatusCreated, "Cloudinit ISO created successfully")
}
