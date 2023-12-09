CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) NOT NULL DEFAULT NOW(),
    name text,
    email text UNIQUE NOT NULL,
    password_hash text,
    activated bool NOT NULL DEFAULT false,
    image text
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp NOT NULL DEFAULT NOW(),
    scope text NOT NULL
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users USING GIN (to_tsvector('simple', email));