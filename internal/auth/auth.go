package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashPass), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return userID, nil
}

func getAuthHeader(headers http.Header, prefix string) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("'Authorization' header does not exist")
	}

	stuff, found := strings.CutPrefix(auth, prefix)
	if !found {

		return "", fmt.Errorf("'%s' not found in 'Authorization' header", prefix)
	}

	return stuff, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	return getAuthHeader(headers, "Bearer ")
}

func GetAPIKey(headers http.Header) (string, error) {
	return getAuthHeader(headers, "ApiKey ")
}

func MakeRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)

	return hex.EncodeToString(b)
}
