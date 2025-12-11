package docs

import _ "embed"

// OpenAPISpec contains the embedded OpenAPI document.
//
//go:embed openapi.yaml
var OpenAPISpec []byte
