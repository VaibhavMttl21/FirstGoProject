-- Create "products" table
CREATE TABLE "products" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL DEFAULT '',
  "price_cents" bigint NOT NULL,
  "currency" character varying(3) NOT NULL,
  "stock_quantity" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
