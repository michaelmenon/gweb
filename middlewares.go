package gweb

import (
	"bytes"
	"fmt"
	"net/http"
)

// WebLogger ... a middleware for looging the request
// internally we use slog
func middlewareLogger(ctx *WebContext) {
	if ctx != nil && ctx.Request != nil {
		logMsg := fmt.Sprintf("%d %s %s %s", ctx.ReplyStatus, ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.WebLog.Info(logMsg)
	}

}

// set the CORS support for the route with the provided headers and methods
func middlewareCorsCustom(ctx *WebContext, headers []string, methods []string) {
	if ctx == nil && ctx.Writer != nil && ctx.Request != nil {
		return
	}
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
	if ctx == nil && ctx.Writer != nil && ctx.Request != nil {
		return
	}
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
