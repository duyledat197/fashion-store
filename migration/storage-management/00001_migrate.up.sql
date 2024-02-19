--  create file table
CREATE TABLE IF NOT EXISTS files(
  "id" serial PRIMARY KEY,
  "mime_type" text,
  "size" bigint,
  "url" text,
  "file_name" text,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

