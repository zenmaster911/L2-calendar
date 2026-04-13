package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/znmaster911/L2-calendar/internal/logger"
	"github.com/znmaster911/L2-calendar/pkg/services"
)

type Handler struct {
	Services *services.Services
	Logger   *logger.Logger
}

func NewHandler(services *services.Services, logg *logger.Logger) *Handler {
	return &Handler{Services: services, Logger: logg}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(h.rLogger)
	// router.Use(middleware.req)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/", h.LogIn)
	})

	// router.Get("/{order_uid}&{date}", h.GetOrder)
	router.Get("/", h.GetEvent)
	router.Route("/events", func(r chi.Router) {
		r.Use(h.userIdentity)
		r.Post("/", h.CreateEvent)
		r.Delete("/{event_id}", h.DeleteEvent)
		r.Patch("/{event_id}", h.UpdateEvent)
	})

	router.Route("/users", func(r chi.Router) {
		r.Post("/", h.NewUser) //
	})

	return router
}
