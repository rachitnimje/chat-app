package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

var secretKey []byte

// InitJWTSecret loads the secretKey from .env file
func InitJWTSecret() {
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	log.Printf("Using HS256 for JWT. Secret Key: %s\n", secretKey)
}

// CreateToken function takes username and return jwt token as string
func CreateToken(username string) (string, error) {
	if len(secretKey) == 0 {
		return "", fmt.Errorf("JWT secret key not initialized")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": username,
		"iss": "yap-up",
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)

	return tokenString, err
}

// VerifyToken take token string, verifies the signature and returns the jwt.Token object
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	return token, err
}

func extractClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	if !token.Valid {
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func ExtractUsername(token *jwt.Token) (string, error) {
	claims, ok := extractClaims(token)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	username, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token")
	}

	return username, nil
}
