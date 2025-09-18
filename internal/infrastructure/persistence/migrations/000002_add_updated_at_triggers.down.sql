-- 000002_add_updated_at_triggers.down.sql

DROP TRIGGER IF EXISTS update_user_updated_at ON "User";
DROP TRIGGER IF EXISTS update_trip_updated_at ON "Trip";
DROP TRIGGER IF EXISTS update_schedule_updated_at ON "Schedule";
DROP TRIGGER IF EXISTS update_sharetoken_updated_at ON "ShareToken";

DROP FUNCTION IF EXISTS update_updated_at_column();