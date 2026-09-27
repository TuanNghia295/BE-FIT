DROP INDEX IF EXISTS idx_refresh_token_family;
DROP INDEX IF EXISTS idx_refresh_token_user_id;

ALTER TABLE refresh_token
    ALTER COLUMN "tokenHash" DROP NOT NULL,
    ALTER COLUMN "tokenFamily" DROP NOT NULL,
    ALTER COLUMN "expiresAt" DROP NOT NULL;