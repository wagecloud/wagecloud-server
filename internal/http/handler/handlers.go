package handler

import (
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

			r.Route("/accounts", func(r chi.Router) {
				r.Get("/{accountID}", func(w http.ResponseWriter, r *http.Request) {})
			})

		})

	})

	return r
}
