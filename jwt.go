package minion

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateJWTToken creates a jwt token with the given secret
func CreateJWTToken(claims map[string]interface{}) (string, error) {
	_, tokenString, err := tokenAuth.Encode(claims)

	return tokenString, err
}

// Authenticator ...
func (c *Context) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		unauthenticated := false

		for _, path := range c.Engine.options.UnauthenticatedRoutes {
			if path == "*" || req.URL.Path == path {
				unauthenticated = true
			}
		}

		if !unauthenticated {
			errResp := struct {
				Code int    `json:"status"`
				Msg  string `json:"message"`
			}{
				http.StatusUnauthorized,
				http.StatusText(http.StatusUnauthorized),
			}

			ctx := req.Context()

			if jwtErr, ok := ctx.Value("jwt.err").(error); ok {
				if jwtErr != nil {
					c.render.JSON(rw, http.StatusUnauthorized, errResp)
					return
				}
			}

			jwtToken, ok := ctx.Value("jwt").(*jwt.Token)
			if !ok || jwtToken == nil || !jwtToken.Valid {
				c.render.JSON(rw, http.StatusUnauthorized, errResp)
				return
			}
		}

		next.ServeHTTP(rw, req)
	})
}
