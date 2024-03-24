package gweb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//write tests

// go test -v -run TestGet
func TestGet(t *testing.T) {

	web := New()
	web.Get("/world", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Hello, world!")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("GET", "/world", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body
	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// go test -v -run TestPost
func TestPost(t *testing.T) {

	web := New()
	web.Post("/save", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Ok")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("POST", "/save", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body
	expected := "Ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// go test -v -run TestWrongMethod
func TestWrongMethod(t *testing.T) {

	web := New()
	web.Post("/save", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Ok")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("Get", "/save", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

}

// go test -v -run TestWrongRoute
func TestWrongRoute(t *testing.T) {

	web := New()
	web.Post("/save", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Ok")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("Post", "/save1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

// go test -v -run TestJSON
func TestJSON(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	u := User{
		Name: "David",
		Age:  30,
	}
	web := New()
	web.Post("/getuser", func(ctx *WebContext) error {

		ctx.Status(200).JSON(u)
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("POST", "/getuser", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body
	var res User

	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("unexpected body")
	}
	if res.Name != "David" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), u)
	}
}

// go test -v -run TestCors
func TestCors(t *testing.T) {

	web := New().WithDefaultCors()
	web.Get("/world", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Hello, world!")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("GET", "/world", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body
	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// go test -v -run TestCustomCors
func TestCustomCors(t *testing.T) {

	web := New().WithCustomCors([]string{"GET"}, []string{"Content-Type", "Authorization"})

	web.Get("/world", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Hello, world!")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("GET", "/world", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body
	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// go test -v -run TestGroupMiddleware
func TestGroupMiddleware(t *testing.T) {

	web := New()
	web.Use(MiddlewareJwt("secret"))
	v1 := web.Group("/v1")

	v1.Get("/world", func(ctx *WebContext) error {

		ctx.Status(200).SendString("Hello, world!")
		return nil
	})
	// Create a new HTTP request to the handler function
	req, err := http.NewRequest("GET", "/v1/world", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	web.WebTest(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

}
