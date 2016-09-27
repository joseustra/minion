package minion

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Middleware middleware type
type Middleware func(http.Handler) http.Handler

// AuthenticatedRoutes authenticate routes using jwt
func AuthenticatedRoutes(jwtToken string, unauthenticatedRoutes []string) HandlerFunc {
	return func(ctx *Context) {
		credentialsOptional := false

		for _, route := range unauthenticatedRoutes {
			if route == "*" {
				credentialsOptional = true
				break
			}

			if route == ctx.Req.URL.RequestURI() {
				credentialsOptional = true
			}
		}

		if !credentialsOptional {
			bearer := ctx.Req.Header.Get("Authorization")
			if len(bearer) > 0 {
				tokenString := bearer[7:len(bearer)]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(jwtToken), nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					context.Set(ctx.Req, "claims", claims)
				} else {
					ctx.Writer.WriteHeader(http.StatusUnauthorized)
					ctx.Writer.Write([]byte("Invalid token"))
					ctx.Writer.WriteHeaderNow()
					ctx.Abort()
					return
				}

				if err != nil {
					ctx.Writer.WriteHeader(http.StatusUnauthorized)
					ctx.Writer.Write([]byte(err.Error()))
					ctx.Writer.WriteHeaderNow()
					ctx.Abort()
					return
				} else if !token.Valid {
					ctx.Writer.WriteHeader(http.StatusUnauthorized)
					ctx.Writer.Write([]byte("invalid authorization token"))
					ctx.Writer.WriteHeaderNow()
					ctx.Abort()
					return
				}
			} else {
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				ctx.Writer.Write([]byte("authorization token required but not present"))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
