package main

import (
	"context"
	"log"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/dto"
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/strategies/sql"
)

func main() {
	db := sql.NewTransactionsPostgresRepository()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	out, err := db.FindOutputsForTransaction(context.Background(), dto.FindTransactionOutput{
		TxID:        "1234",
		IncludeBEEF: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	println(out)
}
