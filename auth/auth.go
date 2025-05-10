package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Claims struct with user_id and expiration as an integer (Unix timestamp)
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var secretKey string
var expirationMinutes int

func init() {
	godotenv.Load(".env")
	secretKey = os.Getenv("SECRET_KEY")
	expMinutes, err := strconv.Atoi(os.Getenv("EXPIRATION_MINUTES"))
	if err != nil {
		panic(fmt.Sprintf("invalid EXPIRATION_MINUTES value: %v", err))
	}
	expirationMinutes = expMinutes
}

// Generate Token function
func GenerateToken(userID string) (string, error) {
	// Calculate expiration time as Unix timestamp (seconds)}
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute).Unix()

	// Create the claims
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set expiration as a NumericDate (internally an integer Unix timestamp)
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}
	// Create the token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	tokenString = "Bearer " + tokenString
	return tokenString, nil
}

// GetUserID extracts the user_id from a JWT token
func GetUserID(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authentication info")
	}
	tokenString := vals[1]
	if tokenString == "" {
		return "", errors.New("no token found")
	}
	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part authentication info")
	}
	// Parse the token with claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// Check for parsing errors
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Return the user_id from claims
	return claims.UserID, nil
}
