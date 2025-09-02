-- Rollback email verification functionality
-- Migration: 20250902100002_add_email_verification

-- Drop email_verifications table and its indexes
DROP INDEX IF EXISTS idx_email_verifications_sent_at;
DROP INDEX IF EXISTS idx_email_verifications_ip;
DROP INDEX IF EXISTS idx_email_verifications_email;
DROP TABLE IF EXISTS email_verifications;

-- Drop email verification index from users table
DROP INDEX IF EXISTS idx_users_email_verify_token;

-- Remove email verification columns from users table
ALTER TABLE users 
    DROP COLUMN IF EXISTS email_verify_expires,
    DROP COLUMN IF EXISTS email_verify_token,
    DROP COLUMN IF EXISTS email_verified;