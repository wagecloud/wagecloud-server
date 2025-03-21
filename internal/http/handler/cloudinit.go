package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/http/response"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/util/file"
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

	cloudinitFile, err := os.Create(req.Userdata.Name + ".iso")
	if err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create temporary file: "+err.Error())
	}
	defer os.Remove(cloudinitFile.Name())

	err = h.service.Cloudinit.CreateCloudinit(cloudinitFile, req.Userdata, req.Metadata)
	if err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to create cloudinit ISO: "+err.Error())
		return
	}

	if err = file.Move(cloudinitFile.Name(), path.Join(config.GetConfig().App.CloudinitDir, cloudinitFile.Name())); err != nil {
		response.FromError(w, http.StatusInternalServerError, "Failed to move cloudinit ISO: "+err.Error())
		return
	}

	response.FromMessage(w, http.StatusCreated, "Cloudinit ISO created successfully")
}
