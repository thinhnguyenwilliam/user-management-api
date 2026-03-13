-- drop trigger
DROP TRIGGER IF EXISTS trg_update_user_updated_at ON users;

-- drop function
DROP FUNCTION IF EXISTS update_user_updated_at();

-- drop table
DROP TABLE IF EXISTS users;

-- optional: drop extension
DROP EXTENSION IF EXISTS "pgcrypto";