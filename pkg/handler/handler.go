package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/znmaster911/L2-calendar/pkg/services"
)

type Handler struct {
	Services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{Services: services}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Get("/{order_uid}&{date}", h.GetOrder)
	router.Get("/", h.GetEvent)
	router.Post("/", h.CreateEvent)
	router.Delete("/", h.Delete) //{event_id}
	router.Patch("/", h.Update)  //{event_id}
	router.Route("/user", func(r chi.Router) {
		r.Post("/", h.NewUser)
	})
	chi.
	return router
}
