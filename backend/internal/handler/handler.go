package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"travelcrm/internal/repository"
	"travelcrm/internal/service"
)

// Handler wires HTTP endpoints to the service layer.
type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// --- Customer / booking endpoints ---

func (h *Handler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.svc.ListRoutes(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, routes)
}

func (h *Handler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	routeID, _ := strconv.ParseInt(q.Get("route_id"), 10, 64)
	origin := strings.TrimSpace(q.Get("origin"))
	destination := strings.TrimSpace(q.Get("destination"))

	schedules, err := h.svc.ListSchedules(r.Context(), routeID, origin, destination)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, schedules)
}

func (h *Handler) ScheduleSeats(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	seatMap, err := h.svc.SeatMap(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, seatMap)
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var in service.BookingInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, err)
		return
	}
	detail, err := h.svc.CreateBooking(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, detail)
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	detail, err := h.svc.GetBooking(r.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (h *Handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	var in service.RescheduleInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, err)
		return
	}
	detail, err := h.svc.Reschedule(r.Context(), code, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

// --- Admin endpoints ---

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.Dashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) AdminCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.svc.ListCustomers(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customers)
}

func (h *Handler) AdminBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.svc.ListBookings(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookings)
}

func (h *Handler) AdminListSchedules(w http.ResponseWriter, r *http.Request) {
	// Admin sees all upcoming schedules (unfiltered).
	schedules, err := h.svc.ListSchedules(r.Context(), 0, "", "")
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, schedules)
}

func (h *Handler) AdminCreateSchedule(w http.ResponseWriter, r *http.Request) {
	var in service.ScheduleInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, err)
		return
	}
	sched, err := h.svc.CreateSchedule(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sched)
}

func (h *Handler) AdminUpdateSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var in service.ScheduleInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, err)
		return
	}
	sched, err := h.svc.UpdateSchedule(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sched)
}

func (h *Handler) AdminVerifyTicket(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	detail, err := h.svc.GetBooking(r.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}
	// A found ticket is valid by construction in this MVP.
	writeJSON(w, http.StatusOK, map[string]any{
		"valid":  true,
		"detail": detail,
	})
}

// idParam parses a required int64 path parameter.
func idParam(r *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	if err != nil || id <= 0 {
		return 0, repository.ErrInvalidInput
	}
	return id, nil
}
