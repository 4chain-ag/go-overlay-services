#!/bin/bash
set -e 
 
oapi-codegen -config ../api/openapi/server/api-cfg.yaml ../api/openapi/server/api.yaml
oapi-codegen -config ../api/openapi/paths/admin/response-models-cfg.yaml ../api/openapi/paths/admin/response-models.yaml

oapi-codegen -config ../api/openapi/paths/non_admin/response-models-cfg.yaml ../api/openapi/paths/non_admin/response-models.yaml
oapi-codegen -config ../api/openapi/paths/non_admin/request-models-cfg.yaml ../api/openapi/paths/non_admin/request-models.yaml

swagger-cli bundle ../api/openapi/server/api.yaml -o docs/api/openapi/bundled.yaml -t yaml
