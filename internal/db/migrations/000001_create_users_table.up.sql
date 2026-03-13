-- user-management-api/internal/db/migrations/000001_create_users_table.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    user_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_fullname VARCHAR(100) NOT NULL,
    user_email VARCHAR(255) UNIQUE NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_age INT CHECK (user_age >= 1 AND user_age <= 150),

    user_status INT NOT NULL DEFAULT 1
        CHECK (user_status IN (1,2,3)),

    user_level INT NOT NULL DEFAULT 3
        CHECK (user_level IN (1,2,3)),

    user_deleted_at TIMESTAMPTZ DEFAULT NULL,
    user_created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN users.user_status IS
'1: active, 2: inactive, 3: banned';

COMMENT ON COLUMN users.user_level IS
'1: admin, 2: moderator, 3: member';

COMMENT ON COLUMN users.user_age IS
'Age of user (1-150)';

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(user_deleted_at);

-- trigger function
CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.user_updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger
CREATE TRIGGER trg_update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_user_updated_at();