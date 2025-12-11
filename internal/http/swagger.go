package apihttp

import (
	"net/http"

	"github.com/hcuri/skool-mvp-app/docs"
)

func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerHTML))
}

func openAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.OpenAPISpec)
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Skool MVP API</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({
        url: '/swagger/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`
