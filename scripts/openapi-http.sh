#!/bin/bash
set -e 
 
oapi-codegen -config ../api/openapi/server/api-cfg.yaml ../api/openapi/server/api.yaml
oapi-codegen -config ../api/openapi/paths/admin/responses-cfg.yaml ../api/openapi/paths/admin/responses.yaml

oapi-codegen -config ../api/openapi/paths/non_admin/responses-cfg.yaml ../api/openapi/paths/non_admin/responses.yaml
oapi-codegen -config ../api/openapi/paths/non_admin/request-bodies-cfg.yaml ../api/openapi/paths/non_admin/request-bodies.yaml

oapi-codegen -config ../api/openapi/errors/responses-cfg.yaml ../api/openapi/errors/responses.yaml

swagger-cli bundle ../api/openapi/server/api.yaml -o ../docs/api/openapi/bundled.yaml -t yaml
