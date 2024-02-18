--  create product table
CREATE TABLE IF NOT EXISTS products(
  "id" bigint PRIMARY KEY,
  "name" text,
  "type" text,
  "image" jsonb,
  "description" text,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

