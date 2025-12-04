package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name     string
	Email    string
	Password string
}

type Node struct {
	UserId     uint    `json:"UserId"`
	Email      string `json:"Email"`
	MacAddress string `json:"MacAddress"`
}

type JwtClaims struct {
	ID           uint   `json:"UserId"`
	Name         string `json:"Name"`
	Email        string `json:"Email"`
	TokenVersion string `json:"TokenVersion"`
	jwt.RegisteredClaims
}
