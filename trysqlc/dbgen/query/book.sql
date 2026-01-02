-- name: InsertBook :one
insert into books (
    id, author, title, description, price, updated_datetime
) values ($1, $2, $3, $4, $5, now()) returning *;


-- name: GetBook :one
select b.*, a.id as author_id, a.created_datetime as author_created_datetime, a.updated_datetime as author_update_datetime 
from books as b join authors as a on b.author = a.id where b.id = $1 limit 1;
