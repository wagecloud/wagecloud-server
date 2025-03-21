package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wagecloud/wagecloud-server/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Domain routes
			r.Route("/domains", func(r chi.Router) {
				r.Post("/", h.CreateDomain)
				r.Post("/{domainID}/start", h.StartDomain)
			})

			// Image routes
			r.Route("/images", func(r chi.Router) {
				r.Post("/", h.CreateImage)
			})

			// Cloudinit routes
			r.Route("/cloudinit", func(r chi.Router) {
				r.Post("/", h.CreateCloudinit)
			})
		})
	})

	return r
}

// respondWithJSON is a helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError is a helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
