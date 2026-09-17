package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"travelcrm/internal/config"
	"travelcrm/internal/middleware"
)

// Routes builds the full HTTP router: middleware, public endpoints, and the
// admin group guarded by Basic Auth.
func Routes(h *Handler, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	// CORS: the web frontend may call the API from a different origin during
	// local development (e.g. a separate static server).
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		// Public customer/booking endpoints.
		r.Get("/routes", h.ListRoutes)
		r.Get("/schedules", h.ListSchedules)
		r.Get("/schedules/{id}/seats", h.ScheduleSeats)
		r.Post("/bookings", h.CreateBooking)
		r.Get("/bookings/{code}", h.GetBooking)
		r.Post("/bookings/{code}/reschedule", h.Reschedule)
		// A booking and its ticket are looked up by the same resolver — a code
		// may be either a booking_code or a ticket_code.
		r.Get("/tickets/{code}", h.GetBooking)
		r.Get("/tickets/{code}/qr", h.GetTicketQR)

		// Admin endpoints, guarded by Basic Auth.
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminBasicAuth(cfg.AdminUser, cfg.AdminPassword))
			r.Get("/admin/dashboard", h.Dashboard)
			r.Get("/admin/customers", h.AdminCustomers)
			r.Get("/admin/bookings", h.AdminBookings)
			r.Get("/admin/schedules", h.AdminListSchedules)
			r.Post("/admin/schedules", h.AdminCreateSchedule)
			r.Put("/admin/schedules/{id}", h.AdminUpdateSchedule)
			r.Get("/admin/tickets/{code}/verify", h.AdminVerifyTicket)
		})
	})

	// Optionally serve a static frontend (single-origin deploy). When
	// STATIC_DIR is set, any non-/api path falls through to the file server,
	// with index.html served for unknown paths (SPA-style).
	if cfg.StaticDir != "" {
		r.Handle("/*", spaFileServer(cfg.StaticDir))
	}

	return r
}

// spaFileServer serves files from dir, falling back to index.html when the
// requested file does not exist so client-side routing works.
func spaFileServer(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	index := filepath.Join(dir, "index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(path); err != nil || info.IsDir() {
			http.ServeFile(w, r, index)
			return
		}
		fs.ServeHTTP(w, r)
	}
}
