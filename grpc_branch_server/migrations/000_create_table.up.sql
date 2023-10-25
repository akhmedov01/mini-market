CREATE TYPE "type_payment" AS ENUM (
  'Card',
  'Cash'
);

CREATE TYPE "type_status" AS ENUM (
  'Success',
  'Cancel'
);

CREATE TYPE "type_staff" AS ENUM (
  'ShopAssistent',
  'Cashier'
);

CREATE TYPE "tarif_type" AS ENUM (
  'Fixed',
  'Percent'
);

CREATE TYPE "type_transaction" AS ENUM (
  'Withdraw',
  'Topup'
);

CREATE TYPE "source_type" AS ENUM (
  'Bonus',
  'Sales'
);

CREATE TABLE IF NOT EXISTS "branches" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "address" varchar,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "sales" (
  "id" uuid PRIMARY KEY,
  "branch_id" uuid,
  "shop_assistent_id" varchar,
  "cashier_id" uuid,
  "payment_type" type_payment,
  "status" type_status,
  "client_name" varchar,
  "price" numeric(10,2),
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "staffs" (
  "id" uuid PRIMARY KEY,
  "branch_id" uuid,
  "tarif_id" uuid,
  "staff_type" type_staff,
  "name" varchar,
  "balance" numeric(10,2),
  "birth_date" date,
  "age" integer,
  "loging" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "staff_tarifs" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "type" tarif_type,
  "amount_for_cash" numeric(10,2),
  "amount_for_card" numeric(10,2),
  "founded_at" date,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "staff_transactions" (
  "id" uuid PRIMARY KEY,
  "sale_id" uuid,
  "staff_id" uuid,
  "transaction_type" type_transaction,
  "source_type" source_type,
  "amount" numeric(10,2),
  "information_about" text,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);  

ALTER TABLE "sales" ADD FOREIGN KEY ("branch_id") REFERENCES "branches" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("shop_assistent_id") REFERENCES "staffs" ("id");

ALTER TABLE "sales" ADD FOREIGN KEY ("cashier_id") REFERENCES "staffs" ("id");

ALTER TABLE "staff_transactions" ADD FOREIGN KEY ("sale_id") REFERENCES "sales" ("id");

ALTER TABLE "staff_transactions" ADD FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id");

ALTER TABLE "staffs" ADD FOREIGN KEY ("tarif_id") REFERENCES "staff_tarifs" ("id");