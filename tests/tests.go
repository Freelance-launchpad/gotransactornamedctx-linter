package main

import (
	"context"
	"database/sql"

	stdlibTransactor "github.com/Thiht/transactor/stdlib"
)

func ok() {

	ctx := context.Background()

	db, _ := sql.Open("pgx", "aaa")

	transactor, _ := stdlibTransactor.NewTransactor(
		db,
		stdlibTransactor.NestedTransactionsSavepoints,
	)
	if err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		return nil
	}); err != nil {
		panic(err)
	}
}

func nok() {

	ctx := context.Background()

	db, _ := sql.Open("pgx", "aaa")

	transactor, _ := stdlibTransactor.NewTransactor(
		db,
		stdlibTransactor.NestedTransactionsSavepoints,
	)
	if err := transactor.WithinTransaction(ctx, func(context.Context) error {
		return nil
	}); err != nil {
		panic(err)
	}
}
