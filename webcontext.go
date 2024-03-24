package gweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strings"
)

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
	wc.ReplyStatus = http.StatusInternalServerError
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
	//wc.Writer.WriteHeader(wc.ReplyStatus)
	err := encoder.Encode(data)
	if err != nil {
		wc.WebLog.Error("sending json", "WebErr", err)
	}
	return nil
}

// SendString ... send the text data
func (wc *WebContext) SendString(data *strings.Reader, contentType ...string) error {

	// Set the Content-Type header to application/json
	if len(contentType) == 0 {
		wc.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		wc.Writer.Header().Set("Content-Type", contentType[0])
	}
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusOK
	}
	//wc.Writer.WriteHeader(wc.ReplyStatus)
	var err error

	_, err = io.Copy(wc.Writer, data)
	if err != nil {
		wc.WebLog.Error("sending string", "WebErr", err)
	}
	return nil
}

// SendString ... send the text data
func (wc *WebContext) SendBytes(data *bytes.Reader) error {

	// Set the Content-Type header to application/json
	wc.Writer.Header().Set("Content-Type", "application/octet-stream")
	if wc.ReplyStatus == 0 {
		wc.ReplyStatus = http.StatusOK
	}
	//wc.Writer.WriteHeader(wc.ReplyStatus)
	var err error

	_, err = io.Copy(wc.Writer, data)
	if err != nil {
		wc.WebLog.Error("sending bytes", "WebErr", err)
	}
	return nil
}

// Render ... render the html data
func (wc *WebContext) RenderString(data *strings.Reader) error {

	return wc.SendString(data, "text/html; charset=utf-8")
}

// RenderFile ... render a file with data using go template
// filePattern ... is a path ot specifc file types like all the htmls in template folder
// filePattern will be template/*.html
// data provide the Data that needs to be passed to the head file
// headFile is the file that is the start of the view for example index.html
// funcMap ... pass any function map that needs to be passed, it is optional
func (wc *WebContext) RenderFiles(filePattern string, data any, headFile string, funcMap template.FuncMap) error {
	templ := template.New("new")
	var err error
	if funcMap != nil {
		templ, err = templ.Funcs(funcMap).ParseGlob(filePattern)
		if err != nil {
			return err
		}
	} else {
		templ, err = templ.ParseGlob(filePattern)
		if err != nil {
			return err
		}
	}

	// Execute the "index.html" template
	err = templ.ExecuteTemplate(wc.Writer, headFile, data)
	if err != nil {
		wc.WebLog.Error("executing template", "WebErr", err)
	}
	return nil
}
