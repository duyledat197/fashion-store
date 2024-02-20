--  create product table
CREATE TABLE IF NOT EXISTS products(
  "id" serial PRIMARY KEY,
  "name" text,
  "type" text,
  "image_urls" text[],
  "description" text,
  "price" float8,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS purchased_products(
  "product_id" serial PRIMARY KEY,
  "user_id" bigint,
  "price" float8,
  "discount" float8,
  "apply_coupon" text,
  "purchase_total" float8,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS purchased_products_user_id_idx ON purchased_products(user_id);

