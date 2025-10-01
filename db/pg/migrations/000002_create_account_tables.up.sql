CREATE TABLE currencies (
    id smallint NOT NULL,
    name char(3),

    PRIMARY KEY (id),
    CHECK (char_length(name) = 3)
);

INSERT INTO currencies VALUES
(1, 'EUR'), (2, 'USD'), (3, 'RUB'), (4, 'RSD');

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
