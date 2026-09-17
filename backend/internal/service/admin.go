package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"travelcrm/internal/model"
	"travelcrm/internal/repository"
)

func (s *Service) Dashboard(ctx context.Context) (repository.DashboardStats, error) {
	return s.repo.Dashboard(ctx)
}

func (s *Service) ListCustomers(ctx context.Context) ([]model.Customer, error) {
	return s.repo.ListCustomers(ctx)
}

func (s *Service) ListBookings(ctx context.Context) ([]model.BookingDetail, error) {
	return s.repo.ListBookingDetails(ctx)
}

// ScheduleInput is the admin request for creating/updating a schedule.
// Origin/destination are used to find-or-create the route on create.
type ScheduleInput struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	DepartureAt time.Time `json:"departure_at"`
	Capacity    int       `json:"capacity"`
}

func (in ScheduleInput) validateForCreate() error {
	if strings.TrimSpace(in.Origin) == "" || strings.TrimSpace(in.Destination) == "" {
		return fmt.Errorf("%w: origin and destination are required", repository.ErrInvalidInput)
	}
	if in.DepartureAt.IsZero() {
		return fmt.Errorf("%w: departure time is required", repository.ErrInvalidInput)
	}
	if in.Capacity <= 0 {
		return fmt.Errorf("%w: capacity must be positive", repository.ErrInvalidInput)
	}
	return nil
}

// CreateSchedule finds-or-creates the route, then inserts the schedule.
func (s *Service) CreateSchedule(ctx context.Context, in ScheduleInput) (model.Schedule, error) {
	if err := in.validateForCreate(); err != nil {
		return model.Schedule{}, err
	}
	routeID, err := s.repo.FindOrCreateRoute(ctx, strings.TrimSpace(in.Origin), strings.TrimSpace(in.Destination))
	if err != nil {
		return model.Schedule{}, err
	}
	return s.repo.CreateSchedule(ctx, routeID, in.DepartureAt, in.Capacity)
}

// UpdateSchedule updates departure time and capacity of an existing schedule.
func (s *Service) UpdateSchedule(ctx context.Context, id int64, in ScheduleInput) (model.Schedule, error) {
	if in.DepartureAt.IsZero() {
		return model.Schedule{}, fmt.Errorf("%w: departure time is required", repository.ErrInvalidInput)
	}
	if in.Capacity <= 0 {
		return model.Schedule{}, fmt.Errorf("%w: capacity must be positive", repository.ErrInvalidInput)
	}
	return s.repo.UpdateSchedule(ctx, id, in.DepartureAt, in.Capacity)
}
