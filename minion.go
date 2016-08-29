package minion

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var l = log.New(os.Stdout, "[minion] ", 0)

// HandlerFunc TODO
type HandlerFunc func(*Context)

// HTMLEngine TODO
type HTMLEngine interface {
	Render(view string, context interface{}, status ...int) error
}

// Engine TODO
type Engine struct {
	*RouterGroup
	router     *httprouter.Router
	allNoRoute []HandlerFunc
	pool       sync.Pool
	options    Options
}

// Version api version
func Version() string {
	return "0.0.1"
}

// Options defines the options to start the API
type Options struct {
	Cors                  []string
	JWTToken              string
	DisableJSONApi        bool
	UnauthenticatedRoutes []string
}

// New returns a new blank Engine instance without any middleware attached.
func New(opts Options) *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		absolutePath: "/",
		engine:       engine,
	}
	engine.options = opts
	engine.router = httprouter.New()
	engine.router.NotFound = http.HandlerFunc(engine.handle404)
	engine.pool.New = func() interface{} {
		ctx := &Context{
			Engine: engine,
			render: render.New(render.Options{
				Layout: "layout",
			}),
		}
		return ctx
	}
	engine.Use(AuthenticatedRoutes(opts.JWTToken, opts.UnauthenticatedRoutes))
	engine.Use(Recovery())
	return engine
}

// Use add middlewares to be used on applicaton
func (c *Engine) Use(middlewares ...HandlerFunc) {
	c.RouterGroup.Use(middlewares...)
	c.allNoRoute = c.combineHandlers(nil)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (c *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(res, req)
}

// Run run the http server.
func (c *Engine) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	l.Printf("Starting server on port [%d]\n", port)
	if err := http.ListenAndServe(addr, WriteLog(c)); err != nil {
		return err
	}
	return nil
}

// RunTLS run the https server.
func (c *Engine) RunTLS(port int, cert string, key string) error {
	addr := fmt.Sprintf(":%d", port)
	l.Printf("Starting server on port [%d]\n", port)
	if err := http.ListenAndServeTLS(addr, cert, key, c); err != nil {
		return err
	}
	return nil
}

func (c *Engine) handle404(w http.ResponseWriter, req *http.Request) {
	ctx := c.createContext(w, req, nil, c.allNoRoute)
	ctx.Writer.WriteHeader(404)
	ctx.Next()
	if !ctx.Writer.Written() {
		if ctx.Writer.Status() == 404 {
			ctx.Writer.Header().Set("Content-Type", "text/html")
			ctx.Writer.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>404 PAGE NOT FOUND</title></head><body style="padding:0;text-align:center;"><div style="padding-top:1em;font-size:2.5em;">404 PAGE NOT FOUND</div><div style="font-size:1em;color:#999;">Powered by Minion</div></body></html>`))
		} else {
			ctx.Writer.WriteHeader(ctx.Writer.Status())
		}
	}
	c.reuseContext(ctx)
}

const (
	// DEV runs the server in development mode
	DEV string = "development"
	// PROD runs the server in production mode
	PROD string = "production"
	// TEST runs the server in test mode
	TEST string = "test"
)

// MinionEnv is the environment that Minion is executing in.
// The MINION_ENV is read on initialization to set this variable.
var MinionEnv = DEV

func init() {
	env := os.Getenv("MINION_ENV")
	if len(env) > 0 {
		MinionEnv = env
	}
}
