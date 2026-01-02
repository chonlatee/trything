-- name: ListAllLink :many
select id, username, amount, link, created_datetime, updated_datetime from link_info;

-- name: ListLinkByUsername :many
select id, username, amount, link, created_datetime, updated_datetime from link_info where username = $1;

