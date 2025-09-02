-- Add email verification functionality
-- Migration: 20250902100002_add_email_verification

-- Add email verification fields to users table
ALTER TABLE users 
    ADD COLUMN email_verified BOOLEAN DEFAULT false,
    ADD COLUMN email_verify_token VARCHAR(255),
    ADD COLUMN email_verify_expires TIMESTAMP;

-- Create index on email_verify_token for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email_verify_token ON users(email_verify_token);

-- Create email_verifications table for tracking verification emails
CREATE TABLE IF NOT EXISTS email_verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for email_verifications table
CREATE INDEX IF NOT EXISTS idx_email_verifications_email ON email_verifications(email);
CREATE INDEX IF NOT EXISTS idx_email_verifications_ip ON email_verifications(ip_address);
CREATE INDEX IF NOT EXISTS idx_email_verifications_sent_at ON email_verifications(sent_at);