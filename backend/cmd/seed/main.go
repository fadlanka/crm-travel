// Command seed resets the database to a known demo state (prompt §15).
// It truncates all tables, then inserts routes, schedules, customers,
// bookings, and tickets with a mix of free and occupied seats.
//
// Run:  go run ./cmd/seed
package main

import (
	"context"
	"log"
	"time"

	"travelcrm/internal/config"
	"travelcrm/internal/db"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("begin: %v", err)
	}
	defer tx.Rollback(ctx)

	// 1. Reset. RESTART IDENTITY resets the id sequences so demo ids are stable.
	if _, err := tx.Exec(ctx, `
		TRUNCATE reschedule_history, tickets, bookings, schedules, routes, customers
		RESTART IDENTITY CASCADE`); err != nil {
		log.Fatalf("truncate: %v", err)
	}

	// 2. Routes.
	routes := []struct{ origin, destination string }{
		{"Jakarta", "Bandung"},
		{"Jakarta", "Semarang"},
		{"Bandung", "Yogyakarta"},
		{"Surabaya", "Malang"},
	}
	routeIDs := make([]int64, len(routes))
	for i, rt := range routes {
		if err := tx.QueryRow(ctx,
			`INSERT INTO routes (origin, destination) VALUES ($1, $2) RETURNING id`,
			rt.origin, rt.destination).Scan(&routeIDs[i]); err != nil {
			log.Fatalf("insert route: %v", err)
		}
	}

	// 3. Schedules. A few today (for the dashboard) and several upcoming.
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	type schedSpec struct {
		routeIdx int
		departAt time.Time
		capacity int
	}
	// Kapasitas 10 kursi per jadwal — cukup untuk demo dan membuat peta kursi ringkas.
	specs := []schedSpec{
		{0, today.Add(8 * time.Hour), 10},  // Jakarta-Bandung today 08:00
		{0, today.Add(15 * time.Hour), 10}, // Jakarta-Bandung today 15:00
		{1, today.Add(9 * time.Hour), 10},  // Jakarta-Semarang today 09:00
		{0, today.AddDate(0, 0, 1).Add(7 * time.Hour), 10},
		{2, today.AddDate(0, 0, 1).Add(10 * time.Hour), 10},
		{3, today.AddDate(0, 0, 2).Add(6 * time.Hour), 10},
		{1, today.AddDate(0, 0, 3).Add(13 * time.Hour), 10},
	}
	schedIDs := make([]int64, len(specs))
	for i, s := range specs {
		if err := tx.QueryRow(ctx,
			`INSERT INTO schedules (route_id, departure_at, capacity) VALUES ($1, $2, $3) RETURNING id`,
			routeIDs[s.routeIdx], s.departAt, s.capacity).Scan(&schedIDs[i]); err != nil {
			log.Fatalf("insert schedule: %v", err)
		}
	}

	// 4. Customers.
	customers := []struct{ name, phone, email string }{
		{"Budi Santoso", "081200000001", "budi@example.com"},
		{"Siti Aminah", "081200000002", "siti@example.com"},
		{"Andi Wijaya", "081200000003", ""},
		{"Dewi Lestari", "081200000004", "dewi@example.com"},
	}
	custIDs := make([]int64, len(customers))
	for i, c := range customers {
		var email *string
		if c.email != "" {
			e := c.email
			email = &e
		}
		if err := tx.QueryRow(ctx,
			`INSERT INTO customers (name, phone, email) VALUES ($1, $2, $3) RETURNING id`,
			c.name, c.phone, email).Scan(&custIDs[i]); err != nil {
			log.Fatalf("insert customer: %v", err)
		}
	}

	// 5. Bookings + tickets. Occupy a spread of seats; leave the rest free.
	//    Codes are fixed so the README can reference them directly.
	type bookingSpec struct {
		custIdx     int
		schedIdx    int
		seat        int
		bookingCode string
		ticketCode  string
	}
	bookings := []bookingSpec{
		{0, 0, 1, "BK-DEMO01", "TK-DEMO0001"},
		{1, 0, 2, "BK-DEMO02", "TK-DEMO0002"},
		{2, 0, 5, "BK-DEMO03", "TK-DEMO0003"},
		{3, 2, 3, "BK-DEMO04", "TK-DEMO0004"},
		{0, 4, 7, "BK-DEMO05", "TK-DEMO0005"},
	}
	for _, b := range bookings {
		var bookingID int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO bookings (customer_id, schedule_id, booking_code, seat_number, status)
			VALUES ($1, $2, $3, $4, 'CONFIRMED') RETURNING id`,
			custIDs[b.custIdx], schedIDs[b.schedIdx], b.bookingCode, b.seat).Scan(&bookingID); err != nil {
			log.Fatalf("insert booking: %v", err)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO tickets (booking_id, ticket_code, status)
			VALUES ($1, $2, 'VALID')`, bookingID, b.ticketCode); err != nil {
			log.Fatalf("insert ticket: %v", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit: %v", err)
	}

	log.Println("seed complete:")
	log.Printf("  %d routes, %d schedules, %d customers, %d bookings",
		len(routes), len(specs), len(customers), len(bookings))
	log.Println("  demo lookup codes: BK-DEMO01 / TK-DEMO0001 (Budi, Jakarta-Bandung)")
}
