-- schema.sql

CREATE TABLE IF NOT EXISTS records
(
    id SERIAL PRIMARY KEY,
    artist VARCHAR (100) NOT NULL,
    album VARCHAR (100) NOT NULL,
    UNIQUE (artist, album)
);

CREATE TABLE IF NOT EXISTS prices
(
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    price NUMERIC(6,2) NOT NULL,
    record_id int NOT NULL REFERENCES records (id),
    UNIQUE (date, record_id)
);
