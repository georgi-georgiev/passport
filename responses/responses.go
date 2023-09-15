package responses

import "math/big"

type IDResp struct {
	ID string `json:"id" example:"1"`
}

type TokenResponse struct {
	TokenType    string `json:"tokenType" example:"Bearer"`
	AccessToken  string `json:"accessToken" example:"token"`
	RefreshToken string `json:"refreshToken" example:"token"`
	ExpiresIn    int    `json:"expiresIn" example:"1687957803"`
}

type ExchangeCodeResponse struct {
	Code string `json:"code" example:"123456"`
}

type CreateUserResponse struct {
	ID           string `json:"id" example:"1"`
	TokenType    string `json:"tokenType" example:"Bearer"`
	AccessToken  string `json:"accessToken" example:"1"`
	RefreshToken string `json:"refreshToken" example:"token"`
	ExpiresIn    int    `json:"expiresIn" example:"1687957803"`
}

type UserResponse struct {
	ID       string   `json:"id" example:"1"`
	Username string   `json:"username" example:"test"`
	Email    string   `json:"email" example:"test@test.com"`
	Role     string   `json:"role" example:"basic"`
	Rights   []string `json:"rights"`
}

type RoleResponse struct {
	Name string `json:"name" example:"admin"`
}

type RightResponse struct {
	Name string `json:"name" example:"default"`
}

type Jwks struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   *big.Int `json:"n"`
	E   int      `json:"e"`
}
