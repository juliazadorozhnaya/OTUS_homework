-- +goose Up
-- SQL in this section is executed when the migration is applied

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Users (
                                     ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    FirstName CHARACTER VARYING(30),
    LastName CHARACTER VARYING(30),
    Email CHARACTER VARYING(30),
    Age INTEGER
    );

CREATE TABLE IF NOT EXISTS Events (
                                      ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    Title CHARACTER VARYING(30),
    Description CHARACTER VARYING(30),
    Beginning TIMESTAMP,
    Finish TIMESTAMP,
    Notification TIMESTAMP,
    UserID UUID,
    FOREIGN KEY (UserID) REFERENCES Users (ID) ON DELETE CASCADE
    );

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

DROP TABLE IF EXISTS Events;

DROP TABLE IF EXISTS Users;

DROP EXTENSION IF EXISTS "uuid-ossp";
