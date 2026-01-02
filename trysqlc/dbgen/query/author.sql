-- name: InsertAuthor :one
insert into authors (
    id, name, bio, updated_datetime
) values ($1, $2, $3, now()) returning *;


-- name: GetAuthor :one
select * from authors where id = $1 limit 1;