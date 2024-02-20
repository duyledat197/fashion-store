CREATE TYPE coupon_type AS ENUM(
  'CouponType_USER',
  'CouponType_PRODUCT',
  'CouponType_LIMITED'
);

CREATE TYPE discount_coupon_type AS ENUM(
  'DiscountType_PERCENT',
  'DiscountType_VALUE'
);

--  create coupon table
CREATE TABLE IF NOT EXISTS coupons(
  "id" serial PRIMARY KEY,
  "code" text UNIQUE,
  "from" timestamptz,
  "to" timestamptz,
  "icon_url" text,
  "description" text,
  "used" bigint,
  "total" bigint,
  "value" float8,
  "image_url" text,
  "discount_type" discount_coupon_type,
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
  "type" coupon_type,
  "created_by" bigint,
  "created_at" timestamptz DEFAULT now(),
  "updated_at" timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS used_coupons_coupon_id_idx ON used_coupons(coupon_id);

CREATE INDEX IF NOT EXISTS used_coupons_user_id_idx ON used_coupons(user_id);

