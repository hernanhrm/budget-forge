CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_users_email UNIQUE NULLS NOT DISTINCT (email)
);

CREATE INDEX idx_users_email ON auth.users (LOWER(email));
CREATE INDEX idx_users_created_at ON auth.users (created_at);
