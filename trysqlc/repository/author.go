package repository

import (
	"context"

	"github.com/chonlatee/trysqlc/dbgen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthorRepo struct {
	db *dbgen.Queries
}

type AuthorRepository interface {
	InsertAuthor(ctx context.Context, param dbgen.InsertAuthorParams) (dbgen.Author, error)
	GetAuthor(ctx context.Context, id pgtype.UUID) (dbgen.Author, error)
}

func (a *AuthorRepo) InsertAuthor(ctx context.Context, param dbgen.InsertAuthorParams) (dbgen.Author, error) {
	return a.db.InsertAuthor(ctx, param)
}

func (a *AuthorRepo) GetAuthor(ctx context.Context, id pgtype.UUID) (dbgen.Author, error) {
	return a.db.GetAuthor(ctx, id)
}

func NewAuthorRepository(conn *pgx.Conn) *AuthorRepo {
	return &AuthorRepo{
		db: dbgen.New(conn),
	}
}
