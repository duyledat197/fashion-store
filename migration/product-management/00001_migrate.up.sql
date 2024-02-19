--  create product table
CREATE TABLE IF NOT EXISTS products(
  "id" serial PRIMARY KEY,
  "sku" text,
  "name" text,
  "type" text,
  "image_urls" text[],
  "description" text,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS products_sku_idx ON products(sku);

