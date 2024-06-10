-- +goose Up
-- SQL in this section is executed when the migration is applied

CREATE SCHEMA IF NOT EXISTS calendar;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS calendar.users (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    FirstName TEXT,
    LastName TEXT,
    Email TEXT,
    Age INTEGER
);

CREATE TABLE IF NOT EXISTS calendar.events (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    Title TEXT,
    Description TEXT,
    Beginning TIMESTAMP(0),
    Finish TIMESTAMP(0),
    Notification TIMESTAMP(0),
    UserID UUID,

    FOREIGN KEY (UserID) REFERENCES calendar.users (ID) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

DROP TABLE IF EXISTS events;

DROP TABLE IF EXISTS users;

DROP SCHEMA IF EXISTS calendar CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
