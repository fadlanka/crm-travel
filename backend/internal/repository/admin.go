package repository

import (
	"context"

	"travelcrm/internal/model"
)

// ListCustomers returns all customers, newest first.
func (r *Repository) ListCustomers(ctx context.Context) ([]model.Customer, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, phone, email, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListBookingDetails returns every booking with its full detail, newest first.
func (r *Repository) ListBookingDetails(ctx context.Context) ([]model.BookingDetail, error) {
	rows, err := r.pool.Query(ctx, bookingDetailSelect+` ORDER BY b.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.BookingDetail
	for rows.Next() {
		d, err := scanBookingDetail(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// DashboardStats is the summary shown on the admin dashboard.
type DashboardStats struct {
	TodaySchedules int `json:"today_schedules"`
	TotalBookings  int `json:"total_bookings"`
	TotalCustomers int `json:"total_customers"`
	TicketsIssued  int `json:"tickets_issued"`
}

// Dashboard returns aggregate counts for the admin dashboard.
func (r *Repository) Dashboard(ctx context.Context) (DashboardStats, error) {
	var s DashboardStats
	err := r.pool.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM schedules WHERE departure_at::date = now()::date),
			(SELECT COUNT(*) FROM bookings WHERE status = 'CONFIRMED'),
			(SELECT COUNT(*) FROM customers),
			(SELECT COUNT(*) FROM tickets WHERE status = 'VALID')`,
	).Scan(&s.TodaySchedules, &s.TotalBookings, &s.TotalCustomers, &s.TicketsIssued)
	return s, err
}
