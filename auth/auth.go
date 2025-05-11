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

// ErrorResponse represents an error response format for API
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
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

// GenerateToken generates a JWT token for the user with the given userID.
// @Summary      Generate JWT Token
// @Description  This function generates a JWT token using the user ID, with an expiration time defined by the EXPIRATION_MINUTES environment variable.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {string}  string  "Bearer Token"
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /auth/token/{user_id} [get]
// @Security     BearerAuth
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

// GetUserID extracts the user ID from the JWT token in the Authorization header.
// @Summary      Extract User ID from Token
// @Description  This function extracts the user ID from the JWT token provided in the Authorization header of the request.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Authorization Token (Bearer)"
// @Success      200            {string}  string  "User ID"
// @Failure      400            {object}  ErrorResponse
// @Failure      401            {object}  ErrorResponse
// @Failure      500            {object}  ErrorResponse
// @Router       /auth/user_id [get]
// @Security     BearerAuth
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
