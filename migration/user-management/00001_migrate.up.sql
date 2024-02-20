-- default role type
CREATE TYPE role_type AS ENUM(
  'SUPER_ADMIN',
  'ADMIN',
  'USER'
);

--  create user table
CREATE TABLE IF NOT EXISTS users(
  "id" serial PRIMARY KEY,
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

CREATE TABLE IF NOT EXISTS login_histories(
  "user_id" bigint,
  "ip" text,
  "user_agent" text,
  "access_token" text,
  "login_at" timestamptz DEFAULT now(),
  "logout_at" timestamptz
);

CREATE INDEX IF NOT EXISTS login_histories_user_id_idx ON login_histories(user_id);

