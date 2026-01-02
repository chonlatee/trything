create table authors (
    id uuid primary key not null,
    name varchar(100) not null,
    bio text,
    created_datetime timestamptz not null default now(),
    updated_datetime timestamptz not null
);