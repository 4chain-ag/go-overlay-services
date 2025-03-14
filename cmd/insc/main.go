package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/transaction/broadcaster"
	"github.com/joho/godotenv"
	"github.com/shruggr/go-block-headers-client/client"
	"github.com/shruggr/goverlay/engine"
	"github.com/shruggr/goverlay/storage/postgres"
	"github.com/shruggr/goverlay/topics"
)

var storage *postgres.PostgresStorage

func init() {
	wd, _ := os.Getwd()
	log.Println("CWD:", wd)
	godotenv.Load(fmt.Sprintf(`%s/../../.env`, wd))
	var err error
	if storage, err = postgres.NewPostgresStorage(context.Background(), os.Getenv("POSTGRES")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var callbackUrl *string
	var callbackToken *string
	if os.Getenv("ARC_CALLBACK") != "" {
		url := os.Getenv("ARC_CALLBACK")
		callbackUrl = &url
	}
	if os.Getenv("ARC_TOKEN") != "" {
		token := os.Getenv("ARC_TOKEN")
		callbackToken = &token
	}
	e := &engine.Engine{
		Managers: map[string]topic.TopicManager{
			"lock":   &topics.LockTopicManager{},
			"bitcom": &topics.BitcomTopicManager{},
			"insc":   &topics.InscriptionTopicManager{},
		},
		LookupServices: map[string]lookup.LookupService{},
		Storage:        storage,
		ChainTracker: &client.HeadersClient{
			Url:    os.Getenv("BLOCK_API"),
			ApiKey: os.Getenv("BLOCK_API_KEY"),
		},
		Broadcaster: &broadcaster.Arc{
			ApiUrl:        "https://arc.taal.com/v1",
			WaitFor:       broadcaster.ACCEPTED_BY_NETWORK,
			CallbackUrl:   callbackUrl,
			CallbackToken: callbackToken,
		},
	}
	log.Println(e)
}
