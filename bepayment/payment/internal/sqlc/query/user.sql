-- name: User_GetUserByUserName :one
select id, username, first_name, last_name, created_datetime, updated_datetime from "user" where username = $1;


-- name: User_GetUserByID :one
select id, username, first_name, last_name, created_datetime, updated_datetime from "user" where id = $1;

-- name: User_GetAllUsers :many
select id, username, first_name, last_name, created_datetime, updated_datetime from "user" offset $1 limit $2;