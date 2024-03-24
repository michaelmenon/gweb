```

web := gweb.New().
		WithLogging().
		WithDefaultReaderWriter("localhost:6379", "1")

	web.Get("/hello", sayHello)
	web.Use(gweb.MiddlewareJwt("secret"))
	v1 := web.Group("/v1")
	v1.Use(getTime)
	v1.Get("/hell", sayHelloV1)
	v1.Post("/user", postUser)

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:

				msg, err := web.MessageController.ReadMessageStream()
				if err == nil {
					fmt.Println(msg)
				}
			}
		}
	}()
	web.Run(":8080")

    func getTime(ctx *gweb.WebContext) error {

	log.Println(time.Now())
	return nil
}
func sayHelloV1(ctx *gweb.WebContext) error {
	ctx.WebLog.Info(ctx.Request.URL.Path)

	ctx.Status(400).SendString("BAD")
	return nil
}
func sayHello(ctx *gweb.WebContext) error {
	ctx.WebLog.Info(ctx.Request.URL.Path)

	ctx.Status(400).SendString("BAD")
	return nil
}

func postUser(ctx *gweb.WebContext) error {
	fmt.Println("id:", ctx.Request.Header.Get("id"))
	usr := new(User)
	ctx.ParseBody(usr)

	fmt.Println(usr)
	ctx.JSON(usr)
	return nil
}

```