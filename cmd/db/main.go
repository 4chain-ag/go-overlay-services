package main

import (
	"context"
	"fmt"
	"log"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/dto"
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine/repositories/strategies/sql"
)

func main() {
	db := sql.NewOutputsPostgresRepository()
	defer db.Close()

	dto, err := db.FindOutput(context.Background(), dto.FindOutputDTO{
		TxID:        "123456",
		OutputIndex: 10,
		Topic:       "example_topic",
		Spent:       false,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dto)
}
