rm .\responses\spec-*.json -ErrorAction SilentlyContinue

## specs-001-lista-specs
$response = "responses/specs-001-lista-specs.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/specs"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response
cat $response | jq

Pause

## specs-002-criar-spec
$request = "requests/specs-002-criar-spec.json"
$response = "responses/specs-002-criar-spec.json"
$method = "POST"
$url = "http://localhost:8080/api/v1/specs"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## specs-003-atualizar-spec
$request = "requests/specs-003-atualizar-spec.json"
$response = "responses/specs-003-atualizar-spec.json"
$method = "PUT"
$url = "http://localhost:8080/api/v1/specs/spec-delegacoes-corporativas-v01"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## specs-004-lista-specs
$response = "responses/specs-004-lista-specs.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/specs"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response
cat $response | jq

Pause

## specs-005-recuperar-spec
$response = "responses/specs-005-recuperar-spec.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/specs/spec-delegacoes-corporativas-v01"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response
cat $response | jq

Pause

## specs-006-deletar-spec
$response = "responses/specs-006-deletar-spec.json"
$method = "DELETE"
$url = "http://localhost:8080/api/v1/specs/spec-delegacoes-corporativas-v01"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response
cat $response | jq

Pause
