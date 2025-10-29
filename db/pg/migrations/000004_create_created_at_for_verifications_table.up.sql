ALTER TABLE verifications
ADD COLUMN created_at timestamptz NOT NULL DEFAULT NOW();

CREATE INDEX idx_verifs_createdat ON verifications(created_at);
