package model

import "time"

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     *string   `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Route struct {
	ID          int64     `json:"id"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Schedule struct {
	ID          int64     `json:"id"`
	RouteID     int64     `json:"route_id"`
	DepartureAt time.Time `json:"departure_at"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Denormalized fields populated when listing schedules with their route.
	Origin         string `json:"origin,omitempty"`
	Destination    string `json:"destination,omitempty"`
	SeatsBooked    int    `json:"seats_booked"`
	SeatsAvailable int    `json:"seats_available"`
}

type Booking struct {
	ID          int64     `json:"id"`
	CustomerID  int64     `json:"customer_id"`
	ScheduleID  int64     `json:"schedule_id"`
	BookingCode string    `json:"booking_code"`
	SeatNumber  int       `json:"seat_number"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Ticket struct {
	ID         int64     `json:"id"`
	BookingID  int64     `json:"booking_id"`
	TicketCode string    `json:"ticket_code"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BookingDetail is the full view returned to customers and admins:
// the booking joined with its customer, schedule, route, and ticket.
type BookingDetail struct {
	Booking     Booking  `json:"booking"`
	Ticket      Ticket   `json:"ticket"`
	Customer    Customer `json:"customer"`
	Schedule    Schedule `json:"schedule"`
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
}
