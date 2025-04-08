package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wagecloud/wagecloud-server/internal/logger"
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
				r.Get("/", h.GetListDomains)
				r.Post("/", h.CreateDomain)
				r.Post("/{domainID}/start", h.StartDomain)
				r.Post("/{domainID}/stop", nil)
				r.Put("/{domainID}", h.UpdateDomain)
				r.Delete("/{domainID}", nil)
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

	// Print routes for debugging
	logger.Log.Info("Registered routes:")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Log.Info(fmt.Sprintf("  %s %s", method, route))
		return nil
	})

	return r
}
