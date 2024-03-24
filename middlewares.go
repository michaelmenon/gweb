package gweb

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	ctx.WebLog.Info(logMsg)
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
