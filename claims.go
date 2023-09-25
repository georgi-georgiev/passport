package passport

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	jwt.StandardClaims
	Role       string   `json:"role"`
	Rights     []string `json:"rights"`
	IsAdmin    bool     `json:"isAdmin"`
	IsVerified bool     `json:"isVerified"`
}

func (c UserClaims) Valid() error {
	return c.StandardClaims.Valid()
}
