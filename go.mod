module github.com/4chain-ag/go-overlay-services

go 1.24.1

require (
	github.com/bsv-blockchain/go-sdk v1.1.22
	github.com/jackc/pgx/v5 v5.7.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/bsv-blockchain/go-sdk => ../go-sdk

replace github.com/shruggr/go-block-headers-client => ../go-block-headers-client
