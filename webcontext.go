package gweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (wc *WebContext) WebLog() zerolog.Logger {
	return log.Logger
}

// get the query parameter
func (wc *WebContext) GetParam(key string) string {

	if wc.query == nil {
		wc.query = wc.Request.URL.Query()
	}
	return wc.query.Get(key)
}

// GetPathValue ... Get the path value
func (wc *WebContext) GetPathValue(key string) string {
	return wc.Request.PathValue(key)
}

// SendStatus ... send the stataus to the user
// is you pass  < 200 status it will be automatically set as 200
func (wc *WebContext) Status(status int) *WebContext {
	if status < 200 {
		status = http.StatusOK
	}
	wc.ReplyStatus = status
	return wc
}

// SendError ... sends the error passed as a response with the ReplyStatus set
func (wc *WebContext) SendError(err error) {
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusInternalServerError
	}
	http.Error(wc.Writer, err.Error(), wc.ReplyStatus)
}

// ParseBody .. parse the request body
func (wc *WebContext) ParseBody(data any) error {
	if data == nil {

		return errors.New(InvalidData)
	}

	decoder := json.NewDecoder(wc.Request.Body)
	return decoder.Decode(data)
}

// JSON ... send the data as a JSON object
func (wc *WebContext) JSON(data any) error {

	if data == nil {

		return errors.New(InvalidData)
	}
	// Set the Content-Type header to application/json
	wc.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(wc.Writer)
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusOK
	}
	wc.Writer.WriteHeader(wc.ReplyStatus)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

// SendString ... send the text data
func (wc *WebContext) SendString(data string, contentType ...string) error {

	// Set the Content-Type header to application/json
	if len(contentType) == 0 {
		wc.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		wc.Writer.Header().Set("Content-Type", contentType[0])
	}
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusOK
	}
	wc.Writer.WriteHeader(wc.ReplyStatus)
	var err error
	var n, total int64

	for {
		n, err = io.Copy(wc.Writer, strings.NewReader(data))

		if err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}
		total += n
		if total == int64(len(data)) {
			return nil
		}
	}

}

// SendString ... send the text data
func (wc *WebContext) SendBytes(data []byte) error {

	// Set the Content-Type header to application/json
	wc.Writer.Header().Set("Content-Type", "application/octet-stream")
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusOK
	}
	wc.Writer.WriteHeader(wc.ReplyStatus)
	var err error
	var n, total int64
	for {
		n, err = io.Copy(wc.Writer, bytes.NewReader(data))

		if err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}
		total += n
		if total == int64(len(data)) {
			return nil
		}
	}

}

// Render ... render the html data
func (wc *WebContext) Render(data string) error {

	return wc.SendString(data, "text/html; charset=utf-8")
}
