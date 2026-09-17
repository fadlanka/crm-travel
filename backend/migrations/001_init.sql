-- Travel CRM — initial schema
-- Small, relational, and easy to explain. See TRAVEL_CRM_CODING_PROMPT.md §7.

CREATE TABLE IF NOT EXISTS customers (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    email       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- One customer profile per phone number (used for find-or-create on booking).
CREATE UNIQUE INDEX IF NOT EXISTS customers_phone_key ON customers (phone);

CREATE TABLE IF NOT EXISTS routes (
    id          BIGSERIAL PRIMARY KEY,
    origin      TEXT NOT NULL,
    destination TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (origin, destination)
);

CREATE TABLE IF NOT EXISTS schedules (
    id           BIGSERIAL PRIMARY KEY,
    route_id     BIGINT NOT NULL REFERENCES routes (id) ON DELETE CASCADE,
    departure_at TIMESTAMPTZ NOT NULL,
    capacity     INT NOT NULL CHECK (capacity > 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS schedules_route_idx ON schedules (route_id);
CREATE INDEX IF NOT EXISTS schedules_departure_idx ON schedules (departure_at);

CREATE TABLE IF NOT EXISTS bookings (
    id           BIGSERIAL PRIMARY KEY,
    customer_id  BIGINT NOT NULL REFERENCES customers (id) ON DELETE RESTRICT,
    schedule_id  BIGINT NOT NULL REFERENCES schedules (id) ON DELETE RESTRICT,
    booking_code TEXT NOT NULL UNIQUE,
    seat_number  INT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'CONFIRMED',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- A seat can only be held once per schedule by an active booking.
CREATE UNIQUE INDEX IF NOT EXISTS bookings_seat_unique
    ON bookings (schedule_id, seat_number)
    WHERE status = 'CONFIRMED';

CREATE INDEX IF NOT EXISTS bookings_customer_idx ON bookings (customer_id);
CREATE INDEX IF NOT EXISTS bookings_schedule_idx ON bookings (schedule_id);

CREATE TABLE IF NOT EXISTS tickets (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL UNIQUE REFERENCES bookings (id) ON DELETE CASCADE,
    ticket_code TEXT NOT NULL UNIQUE,
    status      TEXT NOT NULL DEFAULT 'VALID',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reschedule_history (
    id               BIGSERIAL PRIMARY KEY,
    booking_id       BIGINT NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    old_schedule_id  BIGINT NOT NULL REFERENCES schedules (id),
    old_seat_number  INT NOT NULL,
    new_schedule_id  BIGINT NOT NULL REFERENCES schedules (id),
    new_seat_number  INT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS reschedule_history_booking_idx ON reschedule_history (booking_id);
