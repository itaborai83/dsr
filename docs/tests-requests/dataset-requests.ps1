rm .\responses\datasets-*.json -ErrorAction SilentlyContinue

## datasets-001-criar-spec.json
$request = "requests/datasets-001-criar-spec.json"
$response = "responses/specs-001-lista-specs.json"
$method = "POST"
$url = "http://localhost:8080/api/v1/specs"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## datasets-002-criar-dataset.json
$request = "requests/datasets-002-criar-dataset.json"
$response = "responses/datasets-002-criar-dataset.json"
$method = "POST"
$url = "http://localhost:8080/api/v1/datasets"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## datasets-003-listar-datasets.json
$response = "responses/datasets-003-listar-datasets.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/datasets"
ocurl -X $method $url -o $response
cat $response | jq

Pause

## datasets-004-criar-batch-01.json
$request = "requests/datasets-004-criar-batch-01.json"
$response = "responses/datasets-004-criar-batch-01.json"
$method = "POST"
$url = "http://localhost:8080/api/v1/datasets/dataset-delegacoes-corporativas-v01-2024-06-01/batches"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## datasets-005-criar-batch-02.json
$request = "requests/datasets-005-criar-batch-02.json"
$response = "responses/datasets-005-criar-batch-02.json"
$method = "POST"
$url = "http://localhost:8080/api/v1/datasets/dataset-delegacoes-corporativas-v01-2024-06-01/batches"
$ct = "Content-Type: application/json"
ocurl -X $method $url -H $ct -o $response -d @$request
cat $response | jq

Pause

## datasets-006-listar-datasets.json
$response = "responses/datasets-006-listar-datasets.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/datasets"
ocurl -X $method $url -o $response
cat $response | jq

Pause

## datasets-007-listar-batch-ids.json
$response = "responses/datasets-007-listar-batch-ids.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/batches"
ocurl -X $method $url -o $response
cat $response | jq

Pause

## datasets-008-recuperar-batch.json
$response = "responses/datasets-008-recuperar-batch.json"
$method = "GET"
$url = "http://localhost:8080/api/v1/datasets/dataset-delegacoes-corporativas-v01-2024-06-01/batches/batch-0001"
ocurl -X $method $url -o $response
cat $response | jq

Pause
