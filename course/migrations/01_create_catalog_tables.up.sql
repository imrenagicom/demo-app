CREATE TABLE IF NOT EXISTS courses
(
    id           UUID NOT NULL PRIMARY KEY,
    name         VARCHAR,
    slug         VARCHAR,
    description  TEXT,
    status       INT,
    published_at TIMESTAMP with time zone,
    created_at   TIMESTAMP with time zone default now(),
    updated_at   TIMESTAMP with time zone default now(),
    deleted_at   TIMESTAMP with time zone,
    UNIQUE (slug)
);

CREATE INDEX IF NOT EXISTS idx_courses_deleted_at on courses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_courses_published_at on courses (published_at);
CREATE INDEX IF NOT EXISTS idx_courses_status on courses (status);

CREATE TABLE IF NOT EXISTS course_batches
(
    id              UUID                     NOT NULL PRIMARY KEY,
    course_id       UUID,
    name            VARCHAR,
    max_seats       INT NOT NULL default 0,
    available_seats INT,
    price           DOUBLE PRECISION,
    currency        VARCHAR,
    status          INT,
    start_date      TIMESTAMP with time zone,
    end_date        TIMESTAMP with time zone,
    created_at      TIMESTAMP with time zone default now(),
    updated_at      TIMESTAMP with time zone default now(),
    deleted_at      TIMESTAMP with time zone,
    version         BIGINT                   default 0,
    CONSTRAINT fk_courses_id FOREIGN KEY (course_id) references courses
);

CREATE INDEX IF NOT EXISTS idx_course_batches_created_at on course_batches (created_at);
CREATE INDEX IF NOT EXISTS idx_course_batches_deleted_at on course_batches (deleted_at);
CREATE INDEX IF NOT EXISTS idx_course_batches_course_id on course_batches (course_id);
CREATE INDEX IF NOT EXISTS idx_course_batches_status on course_batches (status);

CREATE TABLE IF NOT EXISTS bookings
(
    id              UUID NOT NULL PRIMARY KEY,
    course_id       UUID,
    course_batch_id UUID,
    price           DOUBLE PRECISION,
    currency        VARCHAR(10),
    status          INT,
    reserved_at     TIMESTAMP with time zone,
    paid_at         TIMESTAMP with time zone,
    created_at      TIMESTAMP with time zone default now(),
    updated_at      TIMESTAMP with time zone default now(),
    deleted_at      TIMESTAMP with time zone,
    version         BIGINT                   default 0,
    CONSTRAINT fk_courses_id FOREIGN KEY (course_id) references courses,
    CONSTRAINT fk_course_batches_id FOREIGN KEY (course_batch_id) references course_batches
);

CREATE INDEX IF NOT EXISTS idx_bookings_deleted_at on bookings (deleted_at);
CREATE INDEX IF NOT EXISTS idx_bookings_course_id on bookings (course_id);
CREATE INDEX IF NOT EXISTS idx_bookings_course_batch_id on bookings (course_batch_id);

ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS invoice_number VARCHAR,
    ADD COLUMN IF NOT EXISTS payment_type VARCHAR;
ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS expired_at TIMESTAMP with time zone;
CREATE INDEX IF NOT EXISTS idx_booking_expires_at on bookings (expired_at);

ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS cust_name VARCHAR NOT NULL default '',
    ADD COLUMN IF NOT EXISTS cust_email VARCHAR NOT NULL default '',
    ADD COLUMN IF NOT EXISTS cust_phone VARCHAR;
