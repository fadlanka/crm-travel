package service

import (
	"context"
	"fmt"

	"travelcrm/internal/model"
	"travelcrm/internal/repository"
)

// RescheduleInput is the validated request for rescheduling a booking.
type RescheduleInput struct {
	NewScheduleID int64 `json:"new_schedule_id"`
	NewSeatNumber int   `json:"new_seat_number"`
}

// Reschedule applies the MVP rule: if the ticket/booking exists, the customer
// may move to any available schedule and seat. (See prompt §2 and §8.)
func (s *Service) Reschedule(ctx context.Context, code string, in RescheduleInput) (model.BookingDetail, error) {
	if in.NewScheduleID <= 0 {
		return model.BookingDetail{}, fmt.Errorf("%w: new schedule is required", repository.ErrInvalidInput)
	}
	if in.NewSeatNumber <= 0 {
		return model.BookingDetail{}, fmt.Errorf("%w: new seat is required", repository.ErrInvalidInput)
	}

	// Rule: ticket must exist to reschedule. GetBooking validates the code and
	// resolves it (booking_code or ticket_code) to the booking.
	current, err := s.GetBooking(ctx, code)
	if err != nil {
		return model.BookingDetail{}, err
	}

	// The new seat must be within the new schedule's capacity.
	sched, err := s.repo.GetSchedule(ctx, in.NewScheduleID)
	if err != nil {
		return model.BookingDetail{}, err
	}
	if in.NewSeatNumber > sched.Capacity {
		return model.BookingDetail{}, fmt.Errorf("%w: seat %d is out of range (capacity %d)",
			repository.ErrInvalidInput, in.NewSeatNumber, sched.Capacity)
	}

	return s.repo.Reschedule(ctx, current.Booking.BookingCode, repository.RescheduleParams{
		NewScheduleID: in.NewScheduleID,
		NewSeatNumber: in.NewSeatNumber,
	})
}
