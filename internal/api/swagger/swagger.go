package swagger

import (
	"context"
	"embed"
	"errors"
	"net/http"
)

// Provider отдаёт Swagger/OpenAPI JSON
type Provider interface {
	SwaggerJSON(ctx context.Context) ([]byte, error)
}

// Router — минимальный интерфейс роутера
type Router interface {
	Get(pattern string, handlerFn http.HandlerFunc)
}

// HTTP — API-компонент, который регистрирует swagger-эндпоинты.
type HTTP interface {
	Register(r Router)
}

type httpAPI struct {
	provider Provider
}

func NewHTTP(provider Provider) HTTP {
	return &httpAPI{provider: provider}
}

func (a *httpAPI) Register(r Router) {
	r.Get("/swagger.json", a.handleSwaggerJSON)
	r.Get("/docs", a.handleDocs)
	r.Get("/docs/", a.handleDocs)
}

func (a *httpAPI) handleSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	b, err := a.provider.SwaggerJSON(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func (a *httpAPI) handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(docsHTML))
}

const docsHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <title>OrderStream API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
    <style>html,body{margin:0;padding:0}#swagger-ui{height:100vh}</style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/swagger.json",
        dom_id: "#swagger-ui"
      });
    </script>
  </body>
</html>`

var swaggerFS embed.FS

// EmbeddedProvider — реализация Provider, которая берёт swagger JSON из embed.
type EmbeddedProvider struct{}

func NewEmbeddedProvider() Provider { return EmbeddedProvider{} }

func (EmbeddedProvider) SwaggerJSON(_ context.Context) ([]byte, error) {
	b, err := swaggerFS.ReadFile("orders.swagger.json")
	if err != nil {
		return nil, errors.New("swagger file not found: run swagger generation")
	}
	return b, nil
}
