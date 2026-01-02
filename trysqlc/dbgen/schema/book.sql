create table books (
    id uuid primary key not null,
    author uuid not null,
    title text not null,
    description text,
    price decimal(10, 2) not null,
    created_datetime timestamptz not null default now(),
    updated_datetime timestamptz not null,
    revoke_datetime timestamptz
);