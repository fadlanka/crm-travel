package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"travelcrm/internal/model"
)

// CreateBookingParams carries the input for a new booking.
type CreateBookingParams struct {
	CustomerName  string
	CustomerPhone string
	CustomerEmail *string
	ScheduleID    int64
	SeatNumber    int
	BookingCode   string
	TicketCode    string
}

// isUniqueViolation reports whether err is a Postgres unique-constraint error.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// CreateBooking runs the full booking flow in one transaction:
// find-or-create the customer, reserve the seat, and issue the ticket.
// A unique partial index enforces seat exclusivity even under races —
// a conflict is translated to ErrSeatTaken.
func (r *Repository) CreateBooking(ctx context.Context, p CreateBookingParams) (model.BookingDetail, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.BookingDetail{}, err
	}
	defer tx.Rollback(ctx)

	// 1. Find or create the customer by phone.
	var customerID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO customers (name, phone, email)
		VALUES ($1, $2, $3)
		ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email, updated_at = now()
		RETURNING id`, p.CustomerName, p.CustomerPhone, p.CustomerEmail).Scan(&customerID)
	if err != nil {
		return model.BookingDetail{}, err
	}

	// 2. Reserve the seat by inserting the booking. The partial unique index
	//    (schedule_id, seat_number) WHERE status='CONFIRMED' does the guarding.
	var bookingID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO bookings (customer_id, schedule_id, booking_code, seat_number, status)
		VALUES ($1, $2, $3, $4, 'CONFIRMED')
		RETURNING id`, customerID, p.ScheduleID, p.BookingCode, p.SeatNumber).Scan(&bookingID)
	if err != nil {
		if isUniqueViolation(err) {
			return model.BookingDetail{}, ErrSeatTaken
		}
		return model.BookingDetail{}, err
	}

	// 3. Issue the ticket.
	_, err = tx.Exec(ctx, `
		INSERT INTO tickets (booking_id, ticket_code, status)
		VALUES ($1, $2, 'VALID')`, bookingID, p.TicketCode)
	if err != nil {
		return model.BookingDetail{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.BookingDetail{}, err
	}
	return r.GetBookingDetailByID(ctx, bookingID)
}

// RescheduleParams carries the input for a reschedule.
type RescheduleParams struct {
	NewScheduleID int64
	NewSeatNumber int
}

// Reschedule moves a booking to a new schedule and seat in one transaction.
// The old seat is freed automatically because the same booking row is moved.
// It records a reschedule_history entry and bumps the ticket's updated_at.
func (r *Repository) Reschedule(ctx context.Context, bookingCode string, p RescheduleParams) (model.BookingDetail, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.BookingDetail{}, err
	}
	defer tx.Rollback(ctx)

	// Load and lock the current booking.
	var bookingID, oldScheduleID int64
	var oldSeat int
	err = tx.QueryRow(ctx, `
		SELECT id, schedule_id, seat_number
		FROM bookings
		WHERE booking_code = $1 AND status = 'CONFIRMED'
		FOR UPDATE`, bookingCode).Scan(&bookingID, &oldScheduleID, &oldSeat)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.BookingDetail{}, ErrNotFound
	}
	if err != nil {
		return model.BookingDetail{}, err
	}

	// Move the booking to the new schedule/seat. The partial unique index
	// rejects a seat that is already taken on the target schedule.
	_, err = tx.Exec(ctx, `
		UPDATE bookings
		SET schedule_id = $2, seat_number = $3, updated_at = now()
		WHERE id = $1`, bookingID, p.NewScheduleID, p.NewSeatNumber)
	if err != nil {
		if isUniqueViolation(err) {
			return model.BookingDetail{}, ErrSeatTaken
		}
		return model.BookingDetail{}, err
	}

	// Record history.
	_, err = tx.Exec(ctx, `
		INSERT INTO reschedule_history
			(booking_id, old_schedule_id, old_seat_number, new_schedule_id, new_seat_number)
		VALUES ($1, $2, $3, $4, $5)`,
		bookingID, oldScheduleID, oldSeat, p.NewScheduleID, p.NewSeatNumber)
	if err != nil {
		return model.BookingDetail{}, err
	}

	// Refresh the ticket (kept valid, timestamp bumped).
	_, err = tx.Exec(ctx, `
		UPDATE tickets SET status = 'VALID', updated_at = now()
		WHERE booking_id = $1`, bookingID)
	if err != nil {
		return model.BookingDetail{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.BookingDetail{}, err
	}
	return r.GetBookingDetailByID(ctx, bookingID)
}

// bookingDetailSelect loads a booking joined with everything a ticket needs.
const bookingDetailSelect = `
	SELECT b.id, b.customer_id, b.schedule_id, b.booking_code, b.seat_number, b.status, b.created_at, b.updated_at,
	       t.id, t.booking_id, t.ticket_code, t.status, t.created_at, t.updated_at,
	       c.id, c.name, c.phone, c.email, c.created_at, c.updated_at,
	       s.id, s.route_id, s.departure_at, s.capacity, s.created_at, s.updated_at,
	       r.origin, r.destination
	FROM bookings b
	JOIN tickets t   ON t.booking_id = b.id
	JOIN customers c ON c.id = b.customer_id
	JOIN schedules s ON s.id = b.schedule_id
	JOIN routes r    ON r.id = s.route_id`

func scanBookingDetail(row pgx.Row) (model.BookingDetail, error) {
	var d model.BookingDetail
	err := row.Scan(
		&d.Booking.ID, &d.Booking.CustomerID, &d.Booking.ScheduleID, &d.Booking.BookingCode,
		&d.Booking.SeatNumber, &d.Booking.Status, &d.Booking.CreatedAt, &d.Booking.UpdatedAt,
		&d.Ticket.ID, &d.Ticket.BookingID, &d.Ticket.TicketCode, &d.Ticket.Status, &d.Ticket.CreatedAt, &d.Ticket.UpdatedAt,
		&d.Customer.ID, &d.Customer.Name, &d.Customer.Phone, &d.Customer.Email, &d.Customer.CreatedAt, &d.Customer.UpdatedAt,
		&d.Schedule.ID, &d.Schedule.RouteID, &d.Schedule.DepartureAt, &d.Schedule.Capacity, &d.Schedule.CreatedAt, &d.Schedule.UpdatedAt,
		&d.Origin, &d.Destination,
	)
	if err != nil {
		return model.BookingDetail{}, err
	}
	d.Schedule.Origin = d.Origin
	d.Schedule.Destination = d.Destination
	return d, nil
}

func (r *Repository) GetBookingDetailByID(ctx context.Context, id int64) (model.BookingDetail, error) {
	row := r.pool.QueryRow(ctx, bookingDetailSelect+` WHERE b.id = $1`, id)
	d, err := scanBookingDetail(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.BookingDetail{}, ErrNotFound
	}
	return d, err
}

// GetBookingDetailByCode looks a booking up by either its booking_code or
// ticket_code — customers may hold either identifier.
func (r *Repository) GetBookingDetailByCode(ctx context.Context, code string) (model.BookingDetail, error) {
	row := r.pool.QueryRow(ctx, bookingDetailSelect+`
		WHERE b.booking_code = $1 OR t.ticket_code = $1`, code)
	d, err := scanBookingDetail(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.BookingDetail{}, ErrNotFound
	}
	return d, err
}
