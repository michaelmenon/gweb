
`go get github.com/michaelmenon/gweb`

**The package is still under development. Do not use in production**

***To initialize call*** `gweb.New()`

    web := gweb.New()

**To enable logging**

    web.WithLogging()

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
         ctx.Status(200).SendString("OK")
         return nil
        }

**Render HTML**

```
func index(ctx *gweb.WebContext) error {
    ctx.WebLog.Info("Index html")

    ctx.Render(<h1>Hello</h1>)
    return nil
}
```

**Adding a handler**

`web.Get("/hello", sayHello)`

`web.Post("/getUser",userHandler)`

Supported HTTP Verbs:
**GET POST PUT DELETE OPTIONS PATCH**

**Adding a middleware**

gweb provides a default middleware for handling JWT authenitcation. You can use it with the following syntax:

`web.Use(gweb.MiddlewareJwt("secret"))`

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

    func sayHelloV1(ctx *gweb.WebContext) error {
    
     //Get path value
    
     ctx.WebLog.Info(ctx.GetPathValue(key))
    
     //Get Query param
    
     ctx.WebLog.Info(ctx.GetParam(key))
    
      
    
     ctx.Status(200).SendString("Hello World")
    
     return nil
    
     }

**Sending JSON**

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
