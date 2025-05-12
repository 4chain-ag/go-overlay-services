package main

import "fmt"

//go:generate go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
//go:generate go tool oapi-codegen --config=../api/openapi/server/api-cfg.yaml         ../api/openapi/server/api.yaml
//go:generate go tool oapi-codegen --config=../api/openapi/paths/admin/responses-cfg.yaml ../api/openapi/paths/admin/responses.yaml
//go:generate go tool oapi-codegen --config=../api/openapi/paths/non_admin/responses-cfg.yaml ../api/openapi/paths/non_admin/responses.yaml
//go:generate go tool oapi-codegen --config=../api/openapi/paths/non_admin/request-bodies-cfg.yaml ../api/openapi/paths/non_admin/request-bodies.yaml

func main() {
	fmt.Println("Overlay HTTP API code generated!")
}
