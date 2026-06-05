CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" decimal(15,2) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" decimal(15,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "origin_account_id" bigint NOT NULL,
  "destination_account_id" bigint NOT NULL,
  "amount" decimal(15,2) NOT NULL CHECK (amount > 0),
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("origin_account_id");

CREATE INDEX ON "transfers" ("destination_account_id");

CREATE INDEX ON "transfers" ("origin_account_id", "destination_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'Can be either positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'Must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "transfers" ADD FOREIGN KEY ("origin_account_id") REFERENCES "accounts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "transfers" ADD FOREIGN KEY ("destination_account_id") REFERENCES "accounts" ("id") DEFERRABLE INITIALLY IMMEDIATE;
