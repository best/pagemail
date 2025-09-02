-- Rollback initial schema
-- Migration: 20250902100001_initial_schema

-- Drop indexes first
DROP INDEX IF EXISTS idx_requests_created_at;
DROP INDEX IF EXISTS idx_requests_status;
DROP INDEX IF EXISTS idx_requests_user_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables (reverse order due to foreign keys)
DROP TABLE IF EXISTS email_configs;
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS users;