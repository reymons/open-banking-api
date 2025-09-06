-- Tables
CREATE TABLE currencies (
    id smallint NOT NULL,
    name char(3),

    CHECK (char_length(name) == 3)
);
INSERT INTO currencies VALUES
(1, 'EUR'), (2, 'USD'), (3, 'RUB'), (4, 'RSD');

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

CREATE TABLE accounts (
    id serial,
    client_id integer NOT NULL,
    number varchar(35) NOT NULL,
    currency_id smallint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    balance money NOT NULL,

    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (curreny_id) REFERENCES currencies(id),
    UNIQUE (number)
);

CREATE TABLE cards (
    id serial,
    number char(16) NOT NULL,
    account_id integer NOT NULL,
    expires_at date NOT NULL,
    cvv char(3),

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE (number, account_id),
    CHECK (char_length(number) == 16 AND char_length(cvv) == 3),
    CHECK (expires_at >= NOW())
);

CREATE TABLE transactions (
    id serial,
    amount money NOT NULL,
    account_id integer NOT NULL,
    beneficiary_id integer,

    FOREIGN KEY (accout_id) REFERENCES accounts(id),
    FOREIGN KEY (beneficiary_id) REFERENCES beneficiaries(id),
    CHECK (amount != 0),
    CHECK (beneficiary_id IS NULL OR amount > 0)
);

CREATE TABLE beneficiaries (
    id serial,
    client_details varchar(255),
    bank_details varchar(255),
);

CREATE TABLE partnership_applications (
    id serial,
    client_id integer NOT NULL,
    text varchar(2048) NOT NULL,

    FOREIGN KEY (client_id) REFERENCES clients(id),
    UNIQUE (client_id)
);

