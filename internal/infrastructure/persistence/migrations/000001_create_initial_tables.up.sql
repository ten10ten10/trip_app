-- 000001_create_initial_tables.up.sql

CREATE EXTENSION IF NOT EXISTS "pg_uuidv7";

CREATE TABLE "User" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT false,
    "verification_token_hash" VARCHAR(255) UNIQUE,
    "verification_token_expires_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "Trip" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "user_id" UUID NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "start_date" DATE NOT NULL,
    "end_date" DATE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "end_date_after_start_date" CHECK ("end_date" >= "start_date")
);

CREATE TABLE "Schedule" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "trip_id" UUID NOT NULL REFERENCES "Trip"("id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "start_date_time" TIMESTAMPTZ NOT NULL,
    "end_date_time" TIMESTAMPTZ NOT NULL,
    "memo" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "end_datetime_after_start_datetime" CHECK ("end_date_time" > "start_date_time")
);

CREATE TABLE "Member" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "trip_id" UUID NOT NULL REFERENCES "Trip"("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL
);

CREATE TABLE "ShareToken" (
    "trip_id" UUID PRIMARY KEY REFERENCES "Trip"("id") ON DELETE CASCADE,
    "token_hash" VARCHAR(255) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);