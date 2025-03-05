package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("auth header empty")
	}
	parts := strings.Fields(authHeader)
	if len(parts) < 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", fmt.Errorf("authorization header format must be 'ApiKey {token}'")
	}

	return parts[1], nil
}
