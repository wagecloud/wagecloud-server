package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/service"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

type Handler struct {
	service *service.Service
}

var (
	validate = validator.New()
	decoder  = schema.NewDecoder()
)

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
			r.Route("/vm", func(r chi.Router) {
				r.Get("/", h.ListVMs)
				r.Get("/{vmID}", h.GetVM)
				r.Post("/", h.CreateVM)
				r.Post("/start/{vmID}", h.StartVM)
				r.Post("/{vmID}/stop", h.StopVM)
				r.Put("/{vmID}", h.UpdateVM)
				r.Delete("/{vmID}", h.DeleteVM)
			})

			r.Route("/account", func(r chi.Router) {
				r.Get("/", h.GetAccount)
				r.Route("/user", func(r chi.Router) {
					r.Post("/login", h.LoginUser)
					r.Post("/register", h.RegisterUser)
				})
			})

			r.Route("/os", func(r chi.Router) {
				r.Get("/", nil)
				r.Get("/{osID}", nil)
				// r.Post("/", h.CreateOS)
				// r.Put("/{osID}", h.UpdateOS)
				// r.Delete("/{osID}", h.DeleteOS)
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

func decodeAndValidate(r *http.Request, v any) error {
	if err := decoder.Decode(v, r.URL.Query()); err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}
