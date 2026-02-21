-- +goose Up
-- Add OAuth2 fields to users table
ALTER TABLE users ADD COLUMN google_id VARCHAR(255);
ALTER TABLE users ADD COLUMN github_id VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar_url TEXT;
ALTER TABLE users ADD COLUMN auth_provider VARCHAR(20) DEFAULT 'email';

-- Create indexes for OAuth2 fields
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_auth_provider ON users(auth_provider);

-- Update existing users to have 'email' as auth provider
UPDATE users SET auth_provider = 'email' WHERE auth_provider IS NULL OR auth_provider = '';

-- +goose Down
-- Remove OAuth2 fields from users table
ALTER TABLE users DROP COLUMN google_id;
ALTER TABLE users DROP COLUMN github_id;
ALTER TABLE users DROP COLUMN avatar_url;
ALTER TABLE users DROP COLUMN auth_provider;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_google_id;
DROP INDEX IF EXISTS idx_users_github_id;
DROP INDEX IF EXISTS idx_users_auth_provider;
