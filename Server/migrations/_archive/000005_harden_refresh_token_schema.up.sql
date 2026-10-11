ALTER TABLE refresh_token
    ALTER COLUMN "tokenHash" SET NOT NULL,
    ALTER COLUMN "tokenFamily" SET NOT NULL,
    ALTER COLUMN "expiresAt" SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_refresh_token_family
    ON refresh_token ("tokenFamily");

CREATE INDEX IF NOT EXISTS idx_refresh_token_user_id
    ON refresh_token ("userId");