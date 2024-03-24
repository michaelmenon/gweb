package gweb

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type InvalidToken struct{}

func (it InvalidToken) Error() string {
	return MsgInvalidToken
}

type ExpiredToken struct{}

func (it ExpiredToken) Error() string {
	return MsgExpiredToken
}

// WebLogger ... a middleware for looging the request
// internally we use slog
func middlewareLogger(ctx *WebContext) {
	logMsg := fmt.Sprintf("%d %s %s %s", ctx.ReplyStatus, ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path)
	log.Info().Msg(logMsg)
}

// JWT ... a jwt middleware to authenitcate the incoing request
// we assume that the Token is set as a Bearer Token under Authorization header key

func MiddlewareJwt(secret string) WebHandler {
	return func(ctx *WebContext) error {
		tokenString := ctx.Request.Header.Get(Authorization)
		if len(tokenString) == 0 {
			return InvalidToken{}
		}
		tok := strings.Split(tokenString, " ")
		if len(tok) < 2 {
			return InvalidToken{}
		}
		token, err := jwt.Parse(tok[1], func(token *jwt.Token) (interface{}, error) {

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secret), nil
		})
		if err != nil {
			return ExpiredToken{}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			for k, v := range claims {

				if value, ok := v.(string); ok {
					ctx.Request.Header.Add(k, value)
				}

			}
		} else {
			return InvalidToken{}
		}
		return nil
	}

}

// set the CORS support for the route with the provided headers and methods
func middlewareCorsCustom(ctx *WebContext, headers []string, methods []string) {

	var allowedHeaders bytes.Buffer
	var allowedMethods bytes.Buffer
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
	if headers != nil && len(headers) > 0 {
		totalHeaders := len(headers)
		for index, v := range headers {
			allowedHeaders.WriteString(v)
			if index < totalHeaders-1 {
				allowedHeaders.WriteString(" ")
			}
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders.String())
	}
	if methods != nil && len(methods) > 0 {
		totalMethods := len(methods)

		for index, v := range methods {
			allowedMethods.WriteString(v)
			if index < totalMethods-1 {
				allowedMethods.WriteString(" ")
			}
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods.String())

	}

	// Check if the request method is OPTIONS (preflight request)
	if ctx.Request.Method == "OPTIONS" {
		// Respond to preflight request
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
}

// set the CORS support for the route
func middlewareCorsDefault(ctx *WebContext) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Check if the request method is OPTIONS (preflight request)
	if ctx.Request.Method == "OPTIONS" {
		// Respond to preflight request
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
}
