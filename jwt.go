package minion

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// CreateJWTToken creates a jwt token with the given secret
func CreateJWTToken(secret string, claims interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.(jwt.Claims))
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetClaims returns the claims
func (c *Context) GetClaims() map[string]interface{} {

	claims := context.Get(c.Req, "claims")
	if claims != nil {
		return claims.(jwt.MapClaims)
	}

	return nil
}
