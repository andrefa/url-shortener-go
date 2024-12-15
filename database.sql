CREATE USER shortener_user WITH PASSWORD 'secure_password';

CREATE DATABASE url_shortener;
GRANT ALL PRIVILEGES ON DATABASE url_shortener TO shortener_user;

CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(6) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);