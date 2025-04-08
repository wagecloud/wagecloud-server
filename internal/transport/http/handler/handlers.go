package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/service"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
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

	// setup 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.FromHTTPError(w, http.StatusNotFound)
	})

	// setup 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.FromHTTPError(w, http.StatusMethodNotAllowed)
	})

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Domain routes
			r.Route("/domain", func(r chi.Router) {
				r.Get("/", h.GetListDomains)
				r.Post("/", h.CreateDomain)
				r.Post("/start/{domainID}", h.StartDomain)
				r.Post("/{domainID}/stop", nil)
				r.Put("/{domainID}", h.UpdateDomain)
				r.Delete("/{domainID}", nil)
			})

			// Image routes
			r.Route("/image", func(r chi.Router) {
				r.Post("/", h.CreateImage)
			})

			// Cloudinit routes
			r.Route("/cloudinit", func(r chi.Router) {
				r.Post("/", h.CreateCloudinit)
			})

			r.Route("/account", func(r chi.Router) {
				r.Get("/", h.GetAccount)
				r.Route("/user", func(r chi.Router) {
					r.Post("/login", h.LoginUser)
					r.Delete("/register", h.RegisterUser)
				})
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
