CREATE TYPE registration_type AS ENUM ('google_auth', 'username_password','linkedin');
CREATE TYPE account_status AS ENUM ('ACTIVE', 'INACTIVE','BLOCKED','LOCKED');
CREATE TYPE password_status AS ENUM ('ACTIVE', 'RESET');


create table users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name varchar,
    last_name varchar,
    password BYTEA,
    pin_attempts integer not null default 0,
    registration_type registration_type not null,
    account_status account_status not null default 'ACTIVE',
    password_status password_status not null default 'ACTIVE',
    created_at timestamp not null default now(),
    update_at timestamp not null default now(),
    deleted_at timestamp
);