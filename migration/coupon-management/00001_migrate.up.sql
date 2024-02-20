CREATE TYPE coupon_type AS ENUM(
  'USER',
  'PRODUCT',
  'LIMITED'
);

CREATE TYPE discount_coupon_type AS ENUM(
  'PERCENT',
  'VALUE',
  'LIMITED'
);

--  create coupon table
CREATE TABLE IF NOT EXISTS coupons(
  "id" bigint PRIMARY KEY,
  "code" text UNIQUE,
  "from" timestamptz,
  "to" timestamptz,
  "rules" jsonb,
  "icon_url" text,
  "description" text,
  "used" bigint,
  "total" bigint,
  "value" float8,
  "discount_coupon_type" discount_coupon_type,
  "coupon_type" coupon_type,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

--  create product coupon table
CREATE TABLE IF NOT EXISTS product_coupons(
  "coupon_id" bigint REFERENCES coupons("id") ON DELETE CASCADE,
  "product_id" bigint,
  "created_by" bigint,
  "used" bigint,
  "total" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now(),
  PRIMARY KEY ("coupon_id", "product_id")
);

CREATE INDEX IF NOT EXISTS product_coupons_coupon_id_idx ON product_coupons(coupon_id);

CREATE INDEX IF NOT EXISTS product_coupons_coupon_id_idx ON product_coupons(product_id);

--  create product user table
CREATE TABLE IF NOT EXISTS user_coupons(
  "coupon_id" bigint REFERENCES coupons("id") ON DELETE CASCADE,
  "user_id" bigint,
  "used" bigint,
  "total" bigint,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now(),
  PRIMARY KEY ("coupon_id", "user_id")
);

CREATE INDEX IF NOT EXISTS user_coupons_coupon_id_idx ON user_coupons(coupon_id);

CREATE INDEX IF NOT EXISTS user_coupons_user_id_idx ON user_coupons(user_id);

--  create used coupon table
CREATE TABLE IF NOT EXISTS used_coupons(
  "coupon_id" bigint REFERENCES coupons("id"),
  "user_id" bigint,
  "type" used_coupon_type,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now(),
  PRIMARY KEY ("coupon_id", "user_id")
);

CREATE INDEX IF NOT EXISTS used_coupons_coupon_id_idx ON used_coupons(coupon_id);

CREATE INDEX IF NOT EXISTS used_coupons_user_id_idx ON used_coupons(user_id);

