package utils

import (
	"errors"
	"os"
	"time"

	"github.com/Blitz-Cloud/smsServer/db"
	"github.com/Blitz-Cloud/smsServer/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *db.User) (string, error) {
	issuer := "https://localhost:3000/api"
	if os.Getenv("env") != "dev" {
		issuer = os.Getenv("domain") + "/api"
	}
	iat := time.Now()
	exp := time.Now().Add(7 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JwtClaims{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		TokenVersion: os.Getenv("jwt_token_version"),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{issuer},
			IssuedAt:  jwt.NewNumericDate(iat),
			NotBefore: jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("jwt_secret")))
	return tokenString, err
}

func ValidateToken(tokenString string) (*types.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("jwt_secret")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return &types.JwtClaims{},errors.New("Invalid JWT signature")
	}

	if claims, ok := token.Claims.(*types.JwtClaims); ok {
		if claims.TokenVersion == os.Getenv("jwt_token_version") {
			return claims, nil
		} else {

			return &types.JwtClaims{}, errors.New("JWT version mismatch\nReceived: " + claims.TokenVersion + "\nAccepts: " + os.Getenv("jwt_token_version"))
		}
	} else {
		return &types.JwtClaims{}, err
	}
}
