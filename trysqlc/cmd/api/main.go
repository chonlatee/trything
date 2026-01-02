package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/chonlatee/trysqlc/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func run() error {

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://pglocal:123456@localhost:35432/postgres")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	log.Println("database connection successfully.")

	// author := repository.NewAuthorRepository(conn)
	// v, err := author.InsertAuthor(ctx, dbgen.InsertAuthorParams{
	// 	ID: pgtype.UUID{
	// 		Bytes: uuid.Must(uuid.NewRandom()),
	// 		Valid: true,
	// 	},
	// 	Name: "JK",
	// 	Bio:  pgtype.Text{String: "Women writer"},
	// })

	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("v: %v\n", v)

	// author := repository.NewAuthorRepository(conn)
	// id, err := uuid.Parse("e889e018-77cc-4c00-a27f-994b119072f5")
	// if err != nil {
	// 	return err
	// }

	// v, err := author.GetAuthor(ctx, pgtype.UUID{Bytes: id, Valid: true})
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("name: %v\n", v.Name)
	// fmt.Printf("bio: %v\n", v.Bio.String)
	// fmt.Printf("created: %v\n", v.CreatedDatetime.Time)
	// fmt.Printf("updated: %v\n", v.UpdatedDatetime.Time)

	id, err := uuid.Parse("83a8cb1c-5258-47d3-b2bc-9de32c831623")
	if err != nil {
		return err
	}

	book := repository.NewBookRepository(conn)

	// var p pgtype.Numeric

	// if err := p.Scan("20.20"); err != nil {
	// 	return err
	// }

	// _, err = book.InsertBook(ctx, dbgen.InsertBookParams{
	// 	ID:          pgtype.UUID{Bytes: uuid.Must(uuid.NewRandom()), Valid: true},
	// 	Author:      pgtype.UUID{Bytes: uuid.Must(uuid.NewRandom()), Valid: true},
	// 	Title:       "Harry potter",
	// 	Description: pgtype.Text{String: "harry potter book 1"},
	// 	Price:       p,
	// })
	// if err != nil {
	// 	return err
	// }

	v, err := book.GetBook(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("row not found")
			return nil
		}
		return err
	}

	price, err := v.Price.Float64Value()
	if err != nil {
		return err
	}

	fmt.Printf("title: %v\n", v.Title)
	fmt.Printf("description: %v\n", v.Description.String)
	fmt.Printf("revoke_datetime: %v\n", v.RevokeDatetime.Valid)
	fmt.Printf("price: %v\n", price)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
