-- default role type
CREATE TYPE role_type AS ENUM(
  'SUPER_ADMIN',
  'ADMIN',
  'USER'
);

--  create user table
CREATE TABLE IF NOT EXISTS users(
  "id" bigint PRIMARY KEY,
  "user_name" text UNIQUE,
  "email" text UNIQUE,
  "password" text NOT NULL,
  "name" text,
  "role" role_type NOT NULL,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS users_user_name_idx ON users(user_name);

CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);

