CREATE TABLE clients (
    id serial,
    role smallint NOT NULL DEFAULT 1, -- Client role
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    birth_date date NOT NULL,
    phone varchar(50) NOT NULL,
    email varchar(50) NOT NULL,
    password varchar(512) NOT NULL,
    is_partner bool DEFAULT FALSE,
    created_at timestamptz NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (email),
    UNIQUE (phone),
    CHECK (birth_date < NOW())
);
