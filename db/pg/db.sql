-- Tables
CREATE TABLE currencies (
    id smallint NOT NULL,
    code char(3),

    PRIMARY KEY (id),
    UNIQUE (code),
    CHECK (char_length(name) = 3)
);

INSERT INTO currencies VALUES
(1, 'EUR'), (2, 'USD'), (3, 'RUB'), (4, 'RSD');

CREATE TABLE verifications (
    target varchar(100) NOT NULL, -- user's email/phone, etc.
    data jsonb NOT NULL, -- any data that needs to be proccessed after validation
    code varchar(10) NOT NULL, -- any verification code from 1 to 10 characters
    expires_at timestamptz NOT NULL,

    UNIQUE (target, code),
    CHECK (char_length(code) > 0),
    CHECK (expires_at > NOW())
);

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

CREATE SEQUENCE account_number_seq START 1000 INCREMENT 1;
CREATE TYPE account_number_type AS ENUM ('active', 'inactive', 'frozen');

CREATE TABLE accounts (
    id serial,
    client_id integer NOT NULL,
    currency_id smallint NOT NULL,
    number varchar(35) NOT NULL DEFAULT lpad(nextval('account_number_seq')::text, 10, '0'),
    balance decimal(19,4) NOT NULL,
    status account_number_type,
    created_at timestamptz NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (currency_id) REFERENCES currencies(id),
    UNIQUE (number)
);

CREATE TABLE cards (
    id serial,
    number char(16) NOT NULL,
    account_id integer NOT NULL,
    expires_at date NOT NULL,
    cvv char(3),

    PRIMARY KEY (id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE (number, account_id),
    CHECK (char_length(number) = 16 AND char_length(cvv) = 3),
    CHECK (expires_at >= NOW())
);

CREATE TABLE transactions (
    id serial,
    amount money NOT NULL,
    account_id integer NOT NULL,
    beneficiary_id integer,

    PRIMARY KEY (id),
    FOREIGN KEY (accout_id) REFERENCES accounts(id),
    FOREIGN KEY (beneficiary_id) REFERENCES beneficiaries(id),
    CHECK (amount != 0),
    CHECK (beneficiary_id IS NULL OR amount > 0)
);

CREATE TABLE beneficiaries (
    id serial,
    client_details varchar(255),
    bank_details varchar(255),

    PRIMARY KEY (id)
);

CREATE TABLE partnership_applications (
    id serial,
    client_id integer NOT NULL,
    text varchar(2048) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    UNIQUE (client_id)
);

