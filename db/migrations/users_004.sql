-- db/migrations/users_nnn.sql

DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    registered_datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved INT NOT NULL,
    human_or_service VARCHAR(7) NOT NULL
);


