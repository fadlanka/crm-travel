package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"travelcrm/internal/model"
)

// scheduleSelect is shared SQL that returns a schedule joined with its route
// plus a live count of confirmed seats. Filters are appended by callers.
const scheduleSelect = `
	SELECT s.id, s.route_id, s.departure_at, s.capacity, s.created_at, s.updated_at,
	       r.origin, r.destination,
	       COALESCE(b.booked, 0) AS seats_booked
	FROM schedules s
	JOIN routes r ON r.id = s.route_id
	LEFT JOIN (
		SELECT schedule_id, COUNT(*) AS booked
		FROM bookings
		WHERE status = 'CONFIRMED'
		GROUP BY schedule_id
	) b ON b.schedule_id = s.id`

func scanSchedule(row pgx.Row) (model.Schedule, error) {
	var s model.Schedule
	err := row.Scan(&s.ID, &s.RouteID, &s.DepartureAt, &s.Capacity, &s.CreatedAt, &s.UpdatedAt,
		&s.Origin, &s.Destination, &s.SeatsBooked)
	if err != nil {
		return model.Schedule{}, err
	}
	s.SeatsAvailable = s.Capacity - s.SeatsBooked
	return s, nil
}

// ListSchedules returns upcoming schedules. Optional filters: routeID (0 = any)
// and origin/destination (empty = any). Only future departures are returned.
func (r *Repository) ListSchedules(ctx context.Context, routeID int64, origin, destination string) ([]model.Schedule, error) {
	query := scheduleSelect + `
		WHERE s.departure_at >= now()
		  AND ($1 = 0 OR s.route_id = $1)
		  AND ($2 = '' OR r.origin = $2)
		  AND ($3 = '' OR r.destination = $3)
		ORDER BY s.departure_at`

	rows, err := r.pool.Query(ctx, query, routeID, origin, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Schedule
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// GetSchedule returns a single schedule with route info and booked count.
func (r *Repository) GetSchedule(ctx context.Context, id int64) (model.Schedule, error) {
	row := r.pool.QueryRow(ctx, scheduleSelect+` WHERE s.id = $1`, id)
	s, err := scanSchedule(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Schedule{}, ErrNotFound
	}
	return s, err
}

// BookedSeats returns the set of seat numbers currently held on a schedule.
func (r *Repository) BookedSeats(ctx context.Context, scheduleID int64) ([]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT seat_number FROM bookings
		WHERE schedule_id = $1 AND status = 'CONFIRMED'
		ORDER BY seat_number`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []int
	for rows.Next() {
		var n int
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		seats = append(seats, n)
	}
	return seats, rows.Err()
}

// CreateSchedule inserts a schedule for an existing route.
func (r *Repository) CreateSchedule(ctx context.Context, routeID int64, departureAt time.Time, capacity int) (model.Schedule, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO schedules (route_id, departure_at, capacity)
		VALUES ($1, $2, $3)
		RETURNING id`, routeID, departureAt, capacity)
	var id int64
	if err := row.Scan(&id); err != nil {
		return model.Schedule{}, err
	}
	return r.GetSchedule(ctx, id)
}

// UpdateSchedule updates departure time and capacity of a schedule.
func (r *Repository) UpdateSchedule(ctx context.Context, id int64, departureAt time.Time, capacity int) (model.Schedule, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE schedules
		SET departure_at = $2, capacity = $3, updated_at = now()
		WHERE id = $1`, id, departureAt, capacity)
	if err != nil {
		return model.Schedule{}, err
	}
	if tag.RowsAffected() == 0 {
		return model.Schedule{}, ErrNotFound
	}
	return r.GetSchedule(ctx, id)
}

// FindOrCreateRoute returns the route id for an origin/destination pair,
// creating it if it does not yet exist.
func (r *Repository) FindOrCreateRoute(ctx context.Context, origin, destination string) (int64, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO routes (origin, destination)
		VALUES ($1, $2)
		ON CONFLICT (origin, destination) DO UPDATE SET updated_at = now()
		RETURNING id`, origin, destination)
	var id int64
	err := row.Scan(&id)
	return id, err
}
