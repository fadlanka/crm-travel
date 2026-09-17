# Travel CRM — Coding Prompt

## 1. Project Context

Build a lightweight Travel CRM / Customer Booking Management System for a technical assignment.

The application is based on a real-world travel booking flow familiar to the developer. The main problem being addressed is that customers often need to contact travel admins manually when they want to reschedule a trip.

The MVP therefore focuses on:
- simple guest booking,
- customer/profile data,
- route and schedule selection,
- seat selection,
- digital ticket generation,
- booking lookup,
- self-service rescheduling,
- simple admin-side customer/booking management,
- ticket verification.

This is intentionally a small system. Do NOT turn it into a large enterprise CRM or full travel platform.

---

## 2. Product Decision

### Target users

1. Travel customers / passengers
2. Travel admin/operator

### Main customer problem

A customer who has already booked a trip may need to change the travel schedule, but traditionally has to contact the travel admin manually.

### Main solution

Allow the customer to manage an existing booking and reschedule it themselves, provided that the ticket exists in the system.

### Reschedule rule for MVP

The business rule is intentionally simple:

```text
Ticket exists = customer is allowed to reschedule
```

Do NOT add payment verification, OTP, account login, departure-time restrictions, cancellation policies, refund calculation, or other complex rules unless explicitly requested later.

When rescheduling:
- customer chooses another available schedule,
- customer chooses an available seat,
- booking information is updated,
- the old seat becomes available again,
- the ticket is updated/regenerated.

---

## 3. Authentication Decision

### Customer

No customer account/login is required.

Reason:
- travel can be an infrequent use case,
- account registration creates unnecessary friction for a lightweight MVP,
- customers can access their booking using their booking/ticket identifier.

### Admin

Use a simple admin authentication mechanism if needed for the admin page.

Keep admin authentication simple. Do not build RBAC, SSO, social login, or complex identity management.

---

## 4. Customer Flow

### Initial booking

```text
Open application
    ↓
Choose origin
    ↓
Choose destination
    ↓
Choose travel date/schedule
    ↓
Choose available seat
    ↓
Enter passenger profile
    ↓
Confirm booking
    ↓
Create booking
    ↓
Generate digital ticket
```

### Manage existing booking

```text
Open Manage Booking
    ↓
Enter ticket / booking code
    ↓
Find ticket
    ↓
Show booking detail
    ↓
Click Reschedule
    ↓
Choose new schedule
    ↓
Choose available seat
    ↓
Confirm reschedule
    ↓
Update booking
    ↓
Generate/update ticket
```

### Ticket verification

```text
Admin opens ticket verification
    ↓
Enter / scan ticket code
    ↓
Find ticket
    ↓
Display ticket information
    ↓
Valid / Invalid
```

---

## 5. Admin Flow

Admin MVP should be intentionally simple.

### Dashboard

Show basic information such as:
- today's schedules,
- total bookings,
- total customers,
- tickets issued.

### Customers

Admin can view customer records.

### Bookings

Admin can view booking information and reschedule history/status if available.

### Schedules

Admin can create/update basic travel schedules, including:
- origin,
- destination,
- departure date/time,
- available seats / capacity.

### Ticket Verification

Admin can enter a ticket code and see whether the ticket exists and its basic booking details.

---

## 6. MVP Features

### Must Have

- [ ] Customer booking without account
- [ ] Origin/destination selection
- [ ] Schedule selection
- [ ] Seat selection
- [ ] Passenger profile input
- [ ] Booking creation
- [ ] Ticket code generation
- [ ] Digital ticket page
- [ ] Manage booking
- [ ] Reschedule booking
- [ ] Release old seat after reschedule
- [ ] Reserve new seat after reschedule
- [ ] Admin dashboard
- [ ] Admin customer list
- [ ] Admin booking list
- [ ] Admin schedule management
- [ ] Ticket verification
- [ ] Seed/demo data
- [ ] Database migration
- [ ] Reset/seed instructions

### Explicitly Out of Scope

- [ ] Online payment gateway
- [ ] Refund processing
- [ ] Cancellation workflow
- [ ] Customer account registration
- [ ] OTP authentication
- [ ] WhatsApp integration
- [ ] Email/SMS integration
- [ ] Loyalty points
- [ ] Promo/voucher system
- [ ] Driver management
- [ ] Vehicle fleet management
- [ ] Accounting
- [ ] Advanced analytics
- [ ] AI features
- [ ] Multi-company / multi-tenant architecture

Do not implement out-of-scope features unless the developer explicitly asks for them.

---

## 7. Suggested Core Data Model

Keep the database small and relational.

Suggested entities:

```text
Customer
- id
- name
- phone
- email (optional)
- created_at
- updated_at
```

```text
Route
- id
- origin
- destination
- created_at
- updated_at
```

```text
Schedule
- id
- route_id
- departure_at
- capacity
- created_at
- updated_at
```

```text
Booking
- id
- customer_id
- schedule_id
- booking_code
- seat_number
- status
- created_at
- updated_at
```

```text
Ticket
- id
- booking_id
- ticket_code
- status
- created_at
- updated_at
```

A separate reschedule-history table is optional. If implemented, keep it simple.

Example:

```text
RescheduleHistory
- id
- booking_id
- old_schedule_id
- old_seat_number
- new_schedule_id
- new_seat_number
- created_at
```

Avoid unnecessary normalization or a complex event-sourcing design.

---

## 8. Important Business Logic

### Booking

When a booking is created:
1. Verify the selected schedule exists.
2. Verify the selected seat is still available.
3. Create/find customer profile.
4. Create booking.
5. Reserve the seat.
6. Generate ticket code.
7. Show digital ticket.

### Reschedule

When a customer reschedules:
1. Find ticket by ticket/booking code.
2. If the ticket does not exist, reject the request.
3. If the ticket exists, allow reschedule.
4. Show schedules that are available for the requested new trip.
5. Show available seats.
6. When confirmed, release the old seat.
7. Reserve the new seat.
8. Update booking schedule and seat.
9. Update/regenerate ticket information.

The MVP intentionally does not implement payment verification because the assignment prototype assumes that a ticket is only generated for a completed/paid booking.

---

## 9. Technical Direction

Preferred stack:

### Frontend

Flutter, preferably configured to run as Flutter Web for the assignment demo.

### Backend

Go.

### Database

PostgreSQL.

### API

REST API.

Use a simple, idiomatic Go structure. Avoid over-engineering.

Suggested structure:

```text
/backend
  /cmd
  /internal
    /handler
    /service
    /repository
    /model
    /middleware
  /migrations
  /seed
```

For frontend:

```text
/frontend
  /lib
    /core
    /models
    /services
    /features
      /booking
      /ticket
      /manage_booking
      /admin
```

The exact structure may be adjusted if there is a simpler approach that keeps the project readable.

---

## 10. Technology Selection Rules

The developer wants to learn Flutter and Go, but the project must remain achievable in approximately one day.

Therefore:
- prioritize simple libraries,
- use standard Go patterns where practical,
- do not introduce frameworks just for the sake of using frameworks,
- do not build custom infrastructure,
- do not introduce microservices,
- do not add Redis unless it is genuinely necessary,
- do not add message queues,
- do not introduce GraphQL,
- do not create a separate authentication service.

For ticket generation:
- a ticket code is mandatory,
- QR code can be used for verification,
- PDF generation is optional and should be skipped if it threatens the one-day scope,
- a responsive digital ticket page is sufficient for MVP.

---

## 11. API Suggestions

Keep endpoints straightforward.

### Customer / Booking

```text
GET    /api/routes
GET    /api/schedules
GET    /api/schedules/:id/seats
POST   /api/bookings
GET    /api/bookings/:bookingCode
POST   /api/bookings/:bookingCode/reschedule
```

### Tickets

```text
GET    /api/tickets/:ticketCode
```

### Admin

```text
GET    /api/admin/customers
GET    /api/admin/bookings
GET    /api/admin/schedules
POST   /api/admin/schedules
PUT    /api/admin/schedules/:id
GET    /api/admin/tickets/:ticketCode/verify
```

Adjust endpoint design when necessary, but keep it RESTful and easy to explain.

---

## 12. UI Requirements

Customer UI should have:

1. Home / trip search
2. Schedule selection
3. Seat selection
4. Passenger form
5. Booking confirmation
6. Digital ticket
7. Manage booking
8. Reschedule

Admin UI should have:

1. Login (simple)
2. Dashboard
3. Customers
4. Bookings
5. Schedules
6. Ticket verification

Keep the visual design clean and professional. Do not spend excessive time on animations or visual effects.

---

## 13. Validation

Keep validation limited to what is necessary for correctness.

Examples:
- origin and destination required,
- schedule required,
- seat required,
- passenger name required,
- phone number required,
- duplicate/occupied seat rejected,
- booking/ticket code must exist for manage booking.

Do not create complicated validation rules that are not required by the MVP.

---

## 14. Error Handling

API should return clear HTTP status codes and readable error messages.

Examples:

```text
400 Bad Request
Invalid booking input
```

```text
404 Not Found
Ticket not found
```

```text
409 Conflict
Selected seat is no longer available
```

The frontend should display a simple, understandable error state.

---

## 15. Seed Data

The repository must include demo data so the reviewer does not see an empty application.

Seed at least:
- several routes,
- several schedules,
- several customers,
- several bookings,
- several tickets,
- a mix of available and occupied seats.

Provide one clear command/process to reset the database to the initial demo state.

Example:

```text
go run ./cmd/seed
```

or an equivalent project command.

The exact command should match the implementation.

---

## 16. Development Approach

Build in small milestones.

### Phase 1 — Project setup
- Flutter project
- Go backend
- PostgreSQL connection
- environment configuration

### Phase 2 — Database
- migrations
- models
- seed/demo data

### Phase 3 — Booking
- routes
- schedules
- seats
- customer profile
- booking creation

### Phase 4 — Ticket
- ticket code
- digital ticket page
- optional QR code

### Phase 5 — Reschedule
- manage booking
- schedule selection
- seat selection
- release old seat
- reserve new seat
- update ticket

### Phase 6 — Admin
- dashboard
- customers
- bookings
- schedules
- ticket verification

### Phase 7 — Verification
- manual testing
- fix critical bugs
- prepare deployment
- README

Do not move to a later phase while an earlier phase is fundamentally broken.

---

## 17. Coding Assistant Instructions

The coding assistant is allowed to generate implementation code, but must preserve the product decisions in this document.

Rules:

1. Do not invent additional business features.
2. Do not silently change business rules.
3. Do not introduce unnecessary dependencies.
4. Prefer readable code over clever code.
5. Keep functions/components focused.
6. Keep database relationships understandable.
7. Explain important implementation decisions before major architectural changes.
8. When a simpler implementation is sufficient, prefer the simpler implementation.
9. Keep the application runnable after each milestone.
10. Do not delete working functionality merely to refactor unless there is a clear reason.

After implementing each milestone, provide:
- files created/changed,
- what was implemented,
- how it works,
- how to run/test it,
- known limitations,
- important code concepts the developer should understand for a technical interview.

---

## 18. Documentation to Produce Later

The coding assistant should NOT create the final assignment documents automatically yet.

The developer will separately prepare:

1. Product/MVP document
2. SRS / functional requirements
3. Architecture document
4. Cost analysis
5. Provider comparison
6. AI usage declaration
7. README

The implementation should provide enough technical information to make those documents easy to write later.

---

## 19. Definition of Done

The MVP is considered done when a reviewer can:

### Customer

```text
Open web app
→ choose route
→ choose schedule
→ choose seat
→ enter passenger data
→ create booking
→ receive ticket
```

and:

```text
Open Manage Booking
→ enter valid ticket/booking code
→ view booking
→ reschedule
→ choose new schedule
→ choose free seat
→ confirm
→ receive updated ticket
```

### Admin

```text
Open admin page
→ see customers/bookings/schedules
→ verify ticket code
→ see whether ticket exists
```

### Repository

Must contain:
- source code,
- migrations,
- seed/demo data,
- setup instructions,
- reset instructions,
- environment variable example,
- clear README.

---

## 20. Important Scope Reminder

This project should remain a **small, explainable system**.

The goal is not maximum features.

The goal is to have every important product and technical decision be understandable and defensible in a technical interview.
