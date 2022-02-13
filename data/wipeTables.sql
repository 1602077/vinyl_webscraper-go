-- wipeTables.sql
-- Drops and re-creates tables to create empty tables for testing.

DROP TABLE price;
DROP TABLE record CASCADE;

CREATE TABLE IF NOT EXISTS record
(
    id SERIAL PRIMARY KEY,
    artist VARCHAR (100) NOT NULL,
    album VARCHAR (100) NOT NULL,
    UNIQUE (artist, album)
);

CREATE TABLE IF NOT EXISTS price
(
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    price DOUBLE PRECISION NOT NULL,
    record_id int NOT NULL REFERENCES record (id),
    UNIQUE (date, record_id)
);
