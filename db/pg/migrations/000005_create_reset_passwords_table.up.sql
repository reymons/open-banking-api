CREATE TABLE reset_password_requests (
    token varchar(255),
    client_id integer NOT NULL,
    expires_at timestamptz NOT NULL,

    PRIMARY KEY (token),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    CHECK (expires_at > NOW())
);
