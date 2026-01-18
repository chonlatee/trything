create table "user" (
    id uuid primary key not null,
    username varchar(100) not null,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    created_datetime timestamptz,
    updated_datetime timestamptz
);