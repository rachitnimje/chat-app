package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Middleware is the middleware function to authenticate user for protected routes.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		// get the authorization header from the request
		authHeader := reader.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Authorization Header: Bearer <token>
		// separate the token from "Bearer"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(writer, fmt.Sprintf("Invalid Authorization header format"), http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// verify the token
		token, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// extract username from the token object
		username, err := ExtractUsername(token)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to extract username from token"), http.StatusUnauthorized)
			return
		}

		// add the username to context, which makes it available to next handlers in the chain
		// the username should be a string (this is checked through type assertion in the GetUsernameFromContext func)
		ctx := context.WithValue(reader.Context(), "username", username)
		reader = reader.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(writer, reader)
	})
}

// GetUsernameFromContext retrieved username from the request context
// this is helper function used by handlers
func GetUsernameFromContext(ctx context.Context) (string, error) {
	username, ok := ctx.Value("username").(string)
	if !ok {
		return "", fmt.Errorf("username not found in request context")
	}
	return username, nil
}
