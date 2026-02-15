-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "price" bigint NOT NULL,
  "stock_quantity" bigint NOT NULL,
  "is_active" boolean NOT NULL,
  "category_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
