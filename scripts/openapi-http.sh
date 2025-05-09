#!/bin/bash
set -e 
 
readonly non_admin_path_dir="../api/openapi/paths/non_admin"
readonly admin_path_dir="../api/openapi/paths/admin"
readonly errors_dir="../api/openapi/errors"
readonly server_dir="../api/openapi/server"
readonly swagger_docs_dir="../docs/api/openapi"
readonly swagger_doc_file="overlay-services.yaml"

oapi-codegen -config $server_dir/api-cfg.yaml $server_dir/api.yaml
oapi-codegen -config $admin_path_dir/responses-cfg.yaml $admin_path_dir/responses.yaml

oapi-codegen -config $non_admin_path_dir/responses-cfg.yaml $non_admin_path_dir/responses.yaml
oapi-codegen -config $non_admin_path_dir/request-bodies-cfg.yaml $non_admin_path_dir/request-bodies.yaml

oapi-codegen -config $errors_dir/responses-cfg.yaml $errors_dir/responses.yaml

swagger-cli bundle $server_dir/api.yaml -o $swagger_docs_dir/$swagger_doc_file -t yaml
