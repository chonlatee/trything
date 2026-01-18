-- name: LinkInfo_ListAllLink :many
select id, username, amount, link, description, created_datetime, updated_datetime from link_info;

-- name: LinkInfo_ListLinkByUsername :many
select id, username, amount, link, description, created_datetime, updated_datetime from link_info where username = $1;

-- name: LinkInfo_InsertLinks :exec
insert into link_info (id, username, amount, link, description, created_datetime, updated_datetime) 
select 
    unnest(sqlc.arg(ids)::uuid[]),
    unnest(sqlc.arg(usernames)::varchar(100)[]),
    unnest(sqlc.arg(amounts)::decimal(10, 2)[]),
    unnest(sqlc.arg(links)::text[]),
    unnest(sqlc.arg(descriptions)::text[]),
    unnest(sqlc.arg(createdDatetimes)::timestamptz[]),
    unnest(sqlc.arg(updatedDatetimes)::timestamptz[]);

