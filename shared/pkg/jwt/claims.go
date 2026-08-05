package jwtutil

import "github.com/golang-jwt/jwt/v5"

type AccesClaims struct {
	jwt.RegisteredClaims
    UserID       string   `json:"uid"`
    Handle       string   `json:"handle"`
    Email        string `json:"email"`
    Roles        []string `json:"roles"`
    TokenVersion int64    `json:"ver"`
}

type RefreshClaims struct {
    jwt.RegisteredClaims
    UserID   string `json:"uid"`
}

type TransactionClaims struct {
    jwt.RegisteredClaims          
    UserID string `json:"uid"`
}