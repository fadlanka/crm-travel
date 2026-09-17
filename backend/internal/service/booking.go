package service

import (
	"context"
	"fmt"
	"strings"

	"travelcrm/internal/model"
	"travelcrm/internal/repository"
)

// BookingInput is the validated request for creating a booking.
type BookingInput struct {
	PassengerName  string  `json:"passenger_name"`
	PassengerPhone string  `json:"passenger_phone"`
	PassengerEmail *string `json:"passenger_email"`
	ScheduleID     int64   `json:"schedule_id"`
	SeatNumber     int     `json:"seat_number"`
}

// validate enforces the minimal correctness rules from the prompt (§13).
func (in BookingInput) validate() error {
	if strings.TrimSpace(in.PassengerName) == "" {
		return fmt.Errorf("%w: passenger name is required", repository.ErrInvalidInput)
	}
	if strings.TrimSpace(in.PassengerPhone) == "" {
		return fmt.Errorf("%w: passenger phone is required", repository.ErrInvalidInput)
	}
	if in.ScheduleID <= 0 {
		return fmt.Errorf("%w: schedule is required", repository.ErrInvalidInput)
	}
	if in.SeatNumber <= 0 {
		return fmt.Errorf("%w: seat is required", repository.ErrInvalidInput)
	}
	return nil
}

// CreateBooking validates input, checks the seat is within capacity and free,
// then delegates the transactional write to the repository.
func (s *Service) CreateBooking(ctx context.Context, in BookingInput) (model.BookingDetail, error) {
	if err := in.validate(); err != nil {
		return model.BookingDetail{}, err
	}

	sched, err := s.repo.GetSchedule(ctx, in.ScheduleID)
	if err != nil {
		return model.BookingDetail{}, err
	}
	if in.SeatNumber > sched.Capacity {
		return model.BookingDetail{}, fmt.Errorf("%w: seat %d is out of range (capacity %d)",
			repository.ErrInvalidInput, in.SeatNumber, sched.Capacity)
	}

	name := strings.TrimSpace(in.PassengerName)
	phone := strings.TrimSpace(in.PassengerPhone)

	return s.repo.CreateBooking(ctx, repository.CreateBookingParams{
		CustomerName:  name,
		CustomerPhone: phone,
		CustomerEmail: normalizeEmail(in.PassengerEmail),
		ScheduleID:    in.ScheduleID,
		SeatNumber:    in.SeatNumber,
		BookingCode:   newBookingCode(),
		TicketCode:    newTicketCode(),
	})
}

// GetBooking looks up a booking by booking_code or ticket_code.
func (s *Service) GetBooking(ctx context.Context, code string) (model.BookingDetail, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return model.BookingDetail{}, fmt.Errorf("%w: code is required", repository.ErrInvalidInput)
	}
	return s.repo.GetBookingDetailByCode(ctx, code)
}

func normalizeEmail(email *string) *string {
	if email == nil {
		return nil
	}
	e := strings.TrimSpace(*email)
	if e == "" {
		return nil
	}
	return &e
}
