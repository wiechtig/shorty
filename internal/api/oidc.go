package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getkin/kin-openapi/openapi3filter"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

func NewAuthenticator(verifier *oidc.IDTokenVerifier) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, verifier, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
// https://github.com/oapi-codegen/oapi-codegen/blob/main/examples/authenticated-api/stdhttp/server/jwt_authenticator.go
func Authenticate(ctx context.Context, verifier *oidc.IDTokenVerifier, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		slog.Error("Invalid security scheme", slog.String("expected", "BearerAuth"), slog.String("actual", input.SecuritySchemeName))
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jwt, err := GetJWTFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		slog.Error("getting token", slog.Any("error", err))
		return fmt.Errorf("getting jwt: %w", err)
	}

	// Parse and verify ID Token payload.
	_, err = verifier.Verify(ctx, jwt)
	if err != nil {
		slog.Error("Invalid token", slog.Any("error", err))
		return fmt.Errorf("validating JWT: %w", err)
	}

	// todo
	//// Extract custom claims
	//var claims struct {
	//	Email    string `json:"email"`
	//	Verified bool   `json:"email_verified"`
	//}
	//if err := idToken.Claims(&claims); err != nil {
	//	// handle error
	//}

	return nil
}

// GetJWTFromRequest extracts a JWT string from an Authorization: Bearer <jwt> header
func GetJWTFromRequest(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")

	// Check for the Authorization header.
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	// We expect a header value of the form "Bearer <token>", with 1 space between
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}
