-- Schema for the authentication service.
-- This file is consumed by sqlc (code generation) AND mounted into the
-- Postgres container's /docker-entrypoint-initdb.d so the table is created
-- automatically on first startup.

CREATE TABLE IF NOT EXISTS users (
    user_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name    TEXT        NOT NULL,
    last_name     TEXT        NOT NULL,
    password      TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,
    phone         TEXT        NOT NULL UNIQUE,
    role          TEXT        NOT NULL,
    token         TEXT,
    refresh_token TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
