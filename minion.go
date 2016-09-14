package minion

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/pressly/chi"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

var l = log.New(os.Stdout, "[minion] ", 0)

// HandlerFunc TODO
type HandlerFunc func(*Context)

// Engine TODO
type Engine struct {
	*RouterGroup
	router     *chi.Mux
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
	engine.router = chi.NewRouter()
	//	engine.router.NotFound = http.HandlerFunc(engine.handle404)
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
	crs := cors.New(cors.Options{
		AllowedOrigins: c.options.Cors,
	})

	addr := fmt.Sprintf(":%d", port)
	l.Printf("Starting server on port [%d]\n", port)
	if err := http.ListenAndServe(addr, crs.Handler(WriteLog(c))); err != nil {
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
