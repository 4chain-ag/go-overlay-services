package main

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server"
)

func main() {
	API := server.NewHTTP(&server.Config{
		AdminBearerToken: "12345678secret!",
		Addr:             "localhost",
		Port:             8080,
	})
	API.ListenAndServe()
}
