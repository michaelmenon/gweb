package gweb

import (
	"net/http"
	"net/url"

	"log/slog"
)

// a web handler type
type WebHandler func(wc *WebContext) error

// the gloabal web instance
var w *Web

// Web struct which handles all the connection
type Web struct {
	//an http server
	httpServer  *http.Server
	middlewares []WebHandler

	router *http.ServeMux
	// a client which implement GwebMessageReaderWriter for message passing and receiving
	MessageController GwebMessageReaderWriter
	//enable Gloabl logging
	logging bool
	//enable cors support for all the routes
	defaultCors bool

	//if we need custom cors headers and methods support
	customHeader []string
	custMethods  []string
	WebLog       *slog.Logger
}

type WebGroup struct {
	router      *http.ServeMux
	pattern     string
	w           *Web
	middlewares []WebHandler
}

// WebContext ... the context for each copnnection
type WebContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	query   url.Values
	//set by the handler
	ReplyStatus int
	WebLog      *slog.Logger
}

// GwebMessage received for this Gweb Service
// Data .. the message for this service
// MessageId is the stream id for this message
type GWebMessage struct {
	Data      string
	MessageId string
}

// the redis stream name where we will get the
const GWebRedisStream = "GwebRedisMessageStream"
const GWebRedisStreamEvent = "GwebRedisStreamEvent"
const GWebRedisStreamGroup = "GwebRedisStreamGroup"

// GwebMessageReaderWriter ...
// implement this interface for any client to act as a message service
// by default we have the redis client
type GwebMessageReaderWriter interface {
	//PostMessage ... post the message to the stream
	PostMessage(WebId string, data string) error
	//ReadMessageStream ... read message from the stream
	//at a time 500 max messages can be read
	ReadMessageStream() ([]GWebMessage, error)
}
