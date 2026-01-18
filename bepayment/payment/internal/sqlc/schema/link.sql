create table link_info (
    id uuid primary key not null,
    amount decimal (10, 2) not null,
    username varchar(100) not null,
    link text not null,
    description text,
    created_datetime timestamptz not null default now(),
    updated_datetime timestamptz not null
)