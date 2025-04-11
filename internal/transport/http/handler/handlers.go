package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
			r.Route("/account", func(r chi.Router) {
				r.Get("/", h.GetAccount)
				r.Route("/user", func(r chi.Router) {
					r.Post("/login", h.LoginUser)
					r.Post("/register", h.RegisterUser)
				})
			})

			r.Route("/vm", func(r chi.Router) {
				r.Get("/", h.ListVMs)
				r.Get("/{vmID}", h.GetVM)
				r.Post("/", h.CreateVM)
				r.Patch("/{vmID}", h.UpdateVM)
				r.Delete("/{vmID}", h.DeleteVM)
				r.Post("/start/{vmID}", h.StartVM)
				r.Post("/stop/{vmID}", h.StopVM)
			})

			r.Route("/os", func(r chi.Router) {
				r.Get("/", h.ListOSs)
				r.Get("/{osID}", h.GetOS)
				r.Post("/", h.CreateOS)
				r.Patch("/{osID}", h.UpdateOS)
				r.Delete("/{osID}", h.DeleteOS)
			})

			r.Route("/network", func(r chi.Router) {
				r.Get("/", h.ListNetworks)
				r.Get("/{networkID}", h.GetNetwork)
				r.Post("/", h.CreateNetwork)
				r.Patch("/{networkID}", h.UpdateNetwork)
				r.Delete("/{networkID}", h.DeleteNetwork)
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

func decodeAndValidate(dst any, src map[string][]string) error {
	if err := decoder.Decode(dst, src); err != nil {
		return err
	}

	if err := validate.Struct(dst); err != nil {
		return err
	}

	return nil
}

func decodeAndValidateJSON(dst any, src io.Reader) error {
	if err := json.NewDecoder(src).Decode(dst); err != nil {
		return err
	}

	if err := validate.Struct(dst); err != nil {
		return err
	}

	return nil
}
