package app

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"secret/log"
)

// Application environments
const (
	Production  string = "Production"
	Testing     string = "Testing"
	Development string = "Development"
)

// App contains application data
type App struct {
	Env            string
	Port           string
	redis          Redis
	Server         *http.Server
	DefaultHandler HandlerFunc
	Summery        map[string]*prometheus.SummaryVec
	Counter        map[string]prometheus.Counter
}

var a App

// New creates App instance and sets default handler function
func New() *App {
	a = App{
		Env:     Development,
		Server:  new(http.Server),
		redis:   NewRedis("localhost:6379", "", 0), // TODO: move redis configuration to .env
		Summery: make(map[string]*prometheus.SummaryVec),
		Counter: make(map[string]prometheus.Counter),
		DefaultHandler: func(c *Ctx) error {
			return c.NotFound()
		},
	}
	return &a
}

// StartHTTP starts web server
func (a *App) StartHTTP(address string) error {
	a.Server.Addr = address
	a.Server.Handler = a
	a.Server.ReadTimeout = 5 * time.Second
	a.Server.WriteTimeout = 10 * time.Second
	return a.Server.ListenAndServe()
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := a.NewCtx(r, w)
	defer func() {
		log.Info("defer")
		if ctx.Route == "" {
			return
		}
		if _, ok := a.Counter[ctx.Route]; ok {
			a.Counter[ctx.Route].Inc()
		}
		if _, ok := a.Summery[ctx.Route]; ok {
			duration := time.Since(start)
			a.Summery[ctx.Route].WithLabelValues("duration").Observe(duration.Seconds())
			size, err := strconv.Atoi(w.Header().Get("Content-Length"))
			if err == nil {
				a.Summery[ctx.Route].WithLabelValues("size").Observe(float64(size))
			}
		}
	}()
	l := log.With("method", r.Method)
	l.With("path", r.URL.Path)
	l.With("query", r.URL.Query())
	l.Info("Request")
	err := a.DefaultHandler(ctx)
	if err != nil {
		ctx.Error(err)
	}
}

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	App   *App
	Req   *http.Request
	Res   http.ResponseWriter
	Path  *Path
	Route string
}

// NewCtx returns pointer to Ctx
func (a *App) NewCtx(req *http.Request, resp http.ResponseWriter) *Ctx {
	ctx := &Ctx{
		App:  a,
		Req:  req,
		Res:  resp,
		Path: NewPath(req.URL),
	}
	return ctx
}

// Success sets 200 http status code and calls Respond
func (c *Ctx) Success(body interface{}) error {
	c.Respond(http.StatusOK, body)
	return nil
}

// Respond sends response
func (c *Ctx) Respond(statusCode int, body interface{}) {
	log.Caller(2).With("status code", statusCode).With("body", body).Info("respond")
	if c.Req.Header.Get("Accept") == "application/xml" {
		c.Res.WriteHeader(statusCode)
		c.Res.Header().Set("Content-Type", "application/xml; charset=UTF-8")
		err := xml.NewEncoder(c.Res).Encode(body)
		if err != nil {
			log.WithError(err).Error("could not encode body")
		}
		return
	}
	// default behavior
	c.Res.WriteHeader(statusCode)
	c.Res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(c.Res).Encode(body)
	if err != nil {
		log.WithError(err).Error("could not encode body")
	}
}

// Call calls handler func
func (c *Ctx) Call(f HandlerFunc) error {
	err := f(c)
	if err != nil {
		c.Error(err)
	}
	return nil
}

// Next calls pathed function with next url segment and increases segment index by 1
func (c *Ctx) Next(f HandlerFunc) error {
	c.Path.Increment()
	return c.Call(f)
}
