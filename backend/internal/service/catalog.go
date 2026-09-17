package service

import (
	"context"

	"travelcrm/internal/model"
)

// SeatMap describes which seats on a schedule are taken vs. free.
type SeatMap struct {
	ScheduleID int64 `json:"schedule_id"`
	Capacity   int   `json:"capacity"`
	Booked     []int `json:"booked"`
	Available  []int `json:"available"`
}

func (s *Service) ListRoutes(ctx context.Context) ([]model.Route, error) {
	return s.repo.ListRoutes(ctx)
}

// ListSchedules returns upcoming schedules, optionally filtered by route or by
// origin/destination names.
func (s *Service) ListSchedules(ctx context.Context, routeID int64, origin, destination string) ([]model.Schedule, error) {
	return s.repo.ListSchedules(ctx, routeID, origin, destination)
}

// SeatMap returns the booked and available seats for a schedule.
func (s *Service) SeatMap(ctx context.Context, scheduleID int64) (SeatMap, error) {
	sched, err := s.repo.GetSchedule(ctx, scheduleID)
	if err != nil {
		return SeatMap{}, err
	}
	booked, err := s.repo.BookedSeats(ctx, scheduleID)
	if err != nil {
		return SeatMap{}, err
	}

	takenSet := make(map[int]bool, len(booked))
	for _, n := range booked {
		takenSet[n] = true
	}
	available := make([]int, 0, sched.Capacity)
	for seat := 1; seat <= sched.Capacity; seat++ {
		if !takenSet[seat] {
			available = append(available, seat)
		}
	}
	if booked == nil {
		booked = []int{}
	}
	return SeatMap{
		ScheduleID: scheduleID,
		Capacity:   sched.Capacity,
		Booked:     booked,
		Available:  available,
	}, nil
}
