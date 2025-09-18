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
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE "Trip" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "userId" UUID NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "startDate" DATE NOT NULL,
    "endDate" DATE NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "end_date_after_start_date" CHECK ("endDate" >= "startDate")
);

CREATE TABLE "Schedule" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "tripId" UUID NOT NULL REFERENCES "Trip"("id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "startDateTime" TIMESTAMPTZ NOT NULL,
    "endDateTime" TIMESTAMPTZ NOT NULL,
    "memo" TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "end_datetime_after_start_datetime" CHECK ("endDateTime" > "startDateTime")
);

CREATE TABLE "Member" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "tripId" UUID NOT NULL REFERENCES "Trip"("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL
);

CREATE TABLE "ShareToken" (
    "tripId" UUID PRIMARY KEY REFERENCES "Trip"("id") ON DELETE CASCADE,
    "token_hash" VARCHAR(255) UNIQUE NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT now()
);