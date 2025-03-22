package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/http/response"
	"github.com/wagecloud/wagecloud-server/internal/model"
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
		response.FromError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cloudinitFile, err := os.Create(path.Join(
		config.GetConfig().App.CloudinitDir,
		fmt.Sprintf("%s_%s.iso", req.Userdata.Name, req.Metadata.InstanceID)),
	)
	if err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create temporary file: "+err.Error())
		return
	}
	// defer os.Remove(cloudinitFile.Name())

	if err = h.service.Cloudinit.CreateCloudinit(cloudinitFile, req.Userdata, req.Metadata); err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create cloudinit ISO: "+err.Error())
		return
	}

	// if err = file.Move(cloudinitFile.Name(), path.Join(config.GetConfig().App.CloudinitDir, cloudinitFile.Name())); err != nil {
	// 	response.FromError(w, http.StatusInternalServerError, "Failed to move cloudinit ISO: "+err.Error())
	// 	return
	// }

	response.FromMessage(w, http.StatusCreated, "Cloudinit ISO created successfully")
}
