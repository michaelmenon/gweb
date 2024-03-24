
`go get github.com/michaelmenon/gweb`

**Requires go 1.22.1 or above**

**To run the unit tests on this package**
``go test -v``

**The package is still under development**

***To initialize call*** `gweb.New()`

    web := gweb.New()

**To enable logging**

    web.WithLogging()

**To enable default CORS**

`web := gweb.New().WithDefaultCors()`

**To Use Custom Cors**

`web := gweb.New().WithCustomCors([]string{"GET"}, []string{"Content-Type", "Authorization"})`

**To initialize with logging**
`web := gweb.New().WithLogging()`

**To use a message stream between two services use the default redis message broker provided**

    web.WithDefaultReaderWriter("localhost:6379", "1")

Listen to messages in a seperate goroutine with the following syntax

    msg, err := web.MessageController.ReadMessageStream()

In a seperate gweb service you can push messages with the followin syntax:

    web.PostMessage("uniqueWebID","data")

**Handler type**

Any function with the following signature is a Handler

    func(ctx *gweb.WebContext) error

For example :

       func sayHello(ctx *gweb.WebContext) error {
         ctx.WebLog.Info(ctx.GetPathValue(key))
         ctx.Status(200).SendString(strings.NewReader("OK"))
         return nil
        }

**Adding a handler**

`web.Get("/hello", sayHello)`

`web.Post("/getUser",userHandler)`

`web.Put("/putUser",userHandler)`

`web.Delete("/deleteUser",deleteUserHandler)`

Supported HTTP Verbs:
**GET POST PUT DELETE OPTIONS PATCH**

**Default Middleware**

1. gweb provides a default middleware for handling JWT authenitcation. You can use it with the following syntax:

`web.Use(gweb.MiddlewareJwt("secret"))`

More are planned

**Custom middleware**

Any function with the following signature can be used as middleware:

     func(ctx *gweb.WebContext) error
    
        func customMiddleware(ctx *gweb.WebContext)error{
         return nil
        }

**Adding a middleware**

    web.Use(customMiddleware)

**Grouping routes**

       v1 := web.Group("/v1")
        //getTime is a custom middleware specifically for routes under v1
        v1.Use(getTime)
        
        v1.Get("/hell", sayHelloV1)
        
        v1.Post("/user", postUser)

**Running the webserver**

      web.Run(":8080")

**Sending a string as response**

```
    func sayHelloV1(ctx *gweb.WebContext) error {
    
     //Get path value
     ctx.WebLog.Info(ctx.GetPathValue(key))
     //Get Query param
     ctx.WebLog.Info(ctx.GetParam(key))

     ctx.Status(200).SendString(strings.NewReader("Hello World")) 
     return nil
    
     }
```

**Sending JSON**

```
    func getUser(ctx *gweb.WebContext) error {
    
        usr := new(User)
        if err:= ctx.ParseBody(usr);err!=nil{
            ctx.WebLog.Error("Body Parsing","parse error",err)
            return err
        }
        ctx.WebLog.Info("Data","user",usr)
        ctx.JSON(usr)
        return nil
    
     }
```

**Render HTML String**

```
func index(ctx *gweb.WebContext) error {
    ctx.WebLog.Info("Index html")

    ctx.RenderString(string.NewReader("<h1>Hello</h1>"))
    return nil
}
```

**Render HTML Files**

```
func index(ctx *gweb.WebContext) error {
    
    funcMap := template.FuncMap{
  "upper": func(s string) string {
   return strings.ToUpper(s)
  }, // Register the custom "upper" function
 }

    //all the html files are in the templates/ folder. The main home page is index.html
    data := "World!"
    ctx.Status(200).RenderFiles("templates/*.html", data, "index.html", funcMap)
    return nil
}
```

**To write unit test check the sample below**

```

func TestGet(t *testing.T) {
    web := New()
    web.Get("/world", func(ctx *WebContext) error {

    ctx.Status(200).SendString(strings.NewReader("Hello, world!))
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

```
