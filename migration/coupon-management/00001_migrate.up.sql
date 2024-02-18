--  create coupon table
CREATE TABLE IF NOT EXISTS coupons(
  "id" bigint PRIMARY KEY,
  "from" timestamptz,
  "to" timestamptz,
  "rules" jsonb,
  "description" text,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

--  create product coupon table
CREATE TABLE IF NOT EXISTS product_coupons(
  "coupon_id" bigint PRIMARY KEY,
  "product_id" bigint,
  "created_by" bigint,
  "amount" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

--  create used coupon table
CREATE TABLE IF NOT EXISTS used_coupons(
  "coupon_id" bigint PRIMARY KEY,
  "user_id" bigint,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

