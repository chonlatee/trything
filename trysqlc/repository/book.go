package repository

import (
	"context"

	"github.com/chonlatee/trysqlc/dbgen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookRepo struct {
	db *dbgen.Queries
}

type BookRepository interface {
	InsertBook(ctx context.Context, param dbgen.InsertBookParams) (dbgen.Book, error)
	GetBook(ctx context.Context, id pgtype.UUID) (dbgen.Book, error)
}

func (b *BookRepo) GetBook(ctx context.Context, id pgtype.UUID) (dbgen.GetBookRow, error) {
	return b.db.GetBook(ctx, id)
}

func (b *BookRepo) InsertBook(ctx context.Context, param dbgen.InsertBookParams) (dbgen.Book, error) {
	return b.db.InsertBook(ctx, param)
}

func NewBookRepository(conn *pgx.Conn) *BookRepo {
	return &BookRepo{
		db: dbgen.New(conn),
	}
}
