package minion

import (
	"net/http"
	"regexp"

	"github.com/goware/jwtauth"
)

// CreateJWTToken creates a jwt token with the given secret
func CreateJWTToken(claims map[string]interface{}) (string, error) {
	_, tokenString, err := tokenAuth.Encode(claims)

	return tokenString, err
}

// Authenticator validates the jwt token and return 401 if not
func (c *Context) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		unauthenticated := false

		for _, path := range c.app.options.UnauthenticatedRoutes {
			re := regexp.MustCompile(path)
			if re.MatchString(req.URL.Path) {
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
			token, _, err := jwtauth.FromContext(req.Context())

			if err != nil {
				c.render.JSON(rw, http.StatusUnauthorized, errResp)
				return
			}

			if token == nil || !token.Valid {
				c.render.JSON(rw, http.StatusUnauthorized, errResp)
				return
			}
		}

		next.ServeHTTP(rw, req)
	})
}
