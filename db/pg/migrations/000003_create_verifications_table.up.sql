CREATE TABLE verifications (
    target varchar(100) NOT NULL, -- user's email/phone, etc.
    data jsonb NOT NULL, -- any data that needs to be proccessed after validation
    code varchar(10) NOT NULL, -- any verification code from 1 to 10 characters
    expires_at timestamptz NOT NULL,

    UNIQUE (target, code),
    CHECK (char_length(code) > 0),
    CHECK (expires_at > NOW())
);
