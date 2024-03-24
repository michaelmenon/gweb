package gweb

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

func New() *Web {
	// initweb will be called only once even from different threads
	sync.OnceFunc(func() {
		router := http.NewServeMux()
		httpServer := &http.Server{
			Addr:           "",
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		w = &Web{
			httpServer:  httpServer,
			middlewares: make([]WebHandler, 0),
			router:      router,
			webLog:      slog.Default(),
		}
	})()
	return w
}
func (w *Web) Logger() *slog.Logger {
	return w.webLog
}

// enable global logging for all the routes
func (w *Web) WithLogging() *Web {
	w.logging = true
	return w
}

// enable default cors
func (w *Web) WithDefaultCors() *Web {
	w.defaultCors = true
	return w
}

// apply CORS with custom headers and methods
func (w *Web) WithCustomCors(headers []string, methods []string) *Web {
	if headers != nil {
		w.customHeader = headers
	}
	if methods != nil {
		w.customHeader = methods
	}
	return w
}

// WithDefaultReaderWriter ... use the default Redis reader writer for message passing between services
func (w *Web) WithDefaultReaderWriter(redisHost string, webId string) *Web {

	rdb, err := newRedisStream(redisHost)
	if err != nil {

		log.Fatal("unable to connect to redis ", err)
		return w
	}
	redisWebStream := &webRedisStream{
		Rdb:   rdb,
		WebId: webId,
	}

	return w.WithMessageReaderWriter(redisWebStream)
}

// WithMessageReaderWriter ... set a messagereaderwriter
// by default it comes with Redis Client
// messageChannel ... the channel to which the messages will be pushed
// use the context Done channel to indicate if we needs to stop receiving the messages
func (w *Web) WithMessageReaderWriter(client GwebMessageReaderWriter) *Web {
	w.MessageController = client

	return w
}

// GetLogger ... call it if you need a build in web logger or use your own
func GetWebLogger() (*slog.Logger, error) {
	if w == nil {
		return nil, errors.New(WebNotInitialized)
	}
	return w.webLog, nil
}

// Run ... create a HTTP server and runs it on the host address provided
// host ... it should be in the "ip:port" format
// returns any error thorwn by the http server
func (w *Web) Run(host string) error {
	w.httpServer.Addr = host
	w.webLog.Info("Running Gweb server on " + host)
	return w.httpServer.ListenAndServe()

}

// addRoutes ... adds the route to the default mux
func (w *Web) addRoutes(pattern string, f WebHandler, wg ...*WebGroup) {

	wc := &WebContext{
		middlewares: make([]WebHandler, len(w.middlewares)),

		WebLog: w.webLog,
	}
	//save the middlewares that needs to be called for this route
	copy(wc.middlewares, w.middlewares)
	handler := func(wr http.ResponseWriter, r *http.Request) {

		wc.Request = r
		wc.Writer = wr
		for _, r := range wc.middlewares {

			e := r(wc)
			if e != nil {
				if errors.Is(e, ExpiredToken{}) || errors.Is(e, InvalidToken{}) {
					http.Error(wr, e.Error(), http.StatusUnauthorized)
				} else {
					http.Error(wr, e.Error(), http.StatusBadRequest)
				}
				return
			}
		}
		if w.defaultCors {
			//write cors headers
			middlewareCorsDefault(wc)
		} else if w.custMethods != nil || w.customHeader != nil {
			middlewareCorsCustom(wc, w.customHeader, w.custMethods)
		}
		err := f(wc)
		if err != nil {
			wc.SendError(err)
		}
		if wc.ReplyStatus == 0 {
			wc.ReplyStatus = http.StatusOK
			wc.Writer.WriteHeader(http.StatusOK)
		}
		if w.logging {
			middlewareLogger(wc)
		}
	}
	if len(wg) > 0 {
		copy(wc.middlewares, wg[0].middlewares)
		wg[0].router.HandleFunc(pattern, handler)

	} else {
		w.router.HandleFunc(pattern, handler)
	}

}

// group the routes
func (w *Web) Group(pattern string) *WebGroup {
	v := WebGroup{
		router:      http.NewServeMux(),
		pattern:     pattern,
		w:           w,
		middlewares: make([]WebHandler, 0),
	}
	if !strings.HasPrefix(pattern, "/") {
		log.Fatal(InvalidPath)
	}

	w.router.Handle("/", v.router)
	return &v
}

// add middlewares using Use
func (w *Web) Use(f WebHandler) {
	w.middlewares = append(w.middlewares, f)
}

// GET ... adds a GET handler
func (w *Web) Get(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodGet+" "+pattern, f)
	return nil
}

// Post ... adds a POST handler
func (w *Web) Post(pattern string, f WebHandler) error {

	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodPost+" "+pattern, f)
	return nil

}

// Delete ... adds a DELETE handler
func (w *Web) Delete(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodDelete+" "+pattern, f)
	return nil
}

// Put ... adds a PUT handler
func (w *Web) Put(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodPut+" "+pattern, f)
	return nil

}
func (w *Web) Options(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodOptions+" "+pattern, f)
	return nil

}
func (w *Web) Patch(pattern string, f WebHandler) error {
	if !strings.HasPrefix(pattern, "/") {
		return errors.New(InvalidPath)
	}
	w.addRoutes(http.MethodPatch+" "+pattern, f)
	return nil
}

//for writing unit test

func (w *Web) WebTest(wr http.ResponseWriter, r *http.Request) {
	w.router.ServeHTTP(wr, r)

}
