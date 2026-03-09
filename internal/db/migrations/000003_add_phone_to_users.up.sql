-- user-management-api/internal/db/migrations/000003_add_phone_to_users.up.sql
ALTER TABLE users
ADD COLUMN phone_number VARCHAR(20) NOT NULL DEFAULT '';