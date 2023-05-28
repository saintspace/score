package server

import (
	"fmt"
	"strings"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

// Placeholder for figuring out JWT validation
func getEmailFromAuthorizationHeader(authHeader string, jwkURL string) (email string, err error) {
	const BearerSchema = "Bearer "
	if authHeader == "" {
		return "", fmt.Errorf("no auth header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, BearerSchema))
	if tokenString == "" {
		return "", fmt.Errorf("no token string")
	}

	jwkSet, err := jwk.Fetch(context.Background(), jwkURL)
	if err != nil {
		return "", fmt.Errorf("error while fetching JWK set => %v", err.Error())
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key ID")
		}
		jwtKey, ok := jwkSet.LookupKeyID(keyID)
		if !ok {
			return nil, fmt.Errorf("no key found with ID: %s", keyID)
		}
		var tokenKey interface{}
		if err := jwtKey.Raw(&tokenKey); err != nil {
			return nil, fmt.Errorf("error while getting raw token key: %s", err)
		}
		return tokenKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("parsed token is not valid => %v", err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("claims are not of type MapClaims")
	}
	// Verify and parse the 'email' claim
	email, ok = claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("email claim is not of type string")
	}
	return email, nil
}
