package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/vm"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/middleware/auth"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetOS(w http.ResponseWriter, r *http.Request) {
}


func (h *Handler) ListOSs(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CreateOS(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateOS(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteOS(w http.ResponseWriter, r *http.Request) {
}


