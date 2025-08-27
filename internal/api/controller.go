package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/jackc/pgx/v5/pgxpool"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"wiechtig.com/shorty/internal/store"
)

//go:generate go tool oapi-codegen -config config.api.yml ../../api/api.yml
//go:generate go tool oapi-codegen -config config.types.yml ../../api/api.yml

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ StrictServerInterface = (*Server)(nil)

type Server struct {
	dbPool  *pgxpool.Pool
	dbStore *store.Queries
}

func New(dbPool *pgxpool.Pool, dbStore *store.Queries) Server {
	return Server{
		dbPool,
		dbStore,
	}
}

type OpenAPIHandlerParams struct {
	Mux      *http.ServeMux
	Server   Server
	Verifier *oidc.IDTokenVerifier
	BaseURL  string
}

func OpenAPIHandler(params OpenAPIHandlerParams) http.Handler {
	mux := params.Mux
	server := params.Server
	verifier := params.Verifier
	baseURL := params.BaseURL

	swagger, err := GetSwagger()
	if err != nil {
		slog.ErrorContext(context.Background(), "Error loading swagger spec", slog.Any("error", err))
		panic(err)
	}
	swagger.Servers = []*openapi3.Server{{URL: baseURL, Description: "Needed to match path"}} // https://github.com/oapi-codegen/oapi-codegen/issues/239#issuecomment-1691380644

	strictHandler := NewStrictHandler(server, nil)

	validator := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(verifier),
			},
		})

	handler := HandlerWithOptions(strictHandler, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: mux,
		Middlewares: []MiddlewareFunc{
			validator,
		},
	})

	return handler
}
