package minion

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

var l = log.New(os.Stdout, "[minion] ", 0)

var tokenAuth *jwtauth.JwtAuth

// HandlerFunc TODO
type HandlerFunc func(*Context)

// Middleware middleware type
type Middleware func(http.Handler) http.Handler

// Engine TODO
type Engine struct {
	*Router
	parent      *Router
	middlewares []Middleware
	pool        sync.Pool
	options     Options
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
	Namespace             string
}

// New returns a new blank Engine instance with no middleware attached
func New(opts Options) *Engine {
	namespace := opts.Namespace
	if len(namespace) == 0 {
		namespace = "/"
	}

	engine := &Engine{}
	engine.Router = &Router{
		namespace: namespace,
		engine:    engine,
		mux:       chi.NewRouter(),
	}
	engine.options = opts
	engine.pool.New = func() interface{} {
		ctx := &Context{
			Engine: engine,
			render: render.New(render.Options{
				Layout: "layout",
			}),
		}
		return ctx
	}

	return engine
}

// Classic returns a new Engine instance with basic middlewares
// Recovery, Logger, CORS and JWT
func Classic(opts Options) *Engine {
	engine := New(opts)
	crs := cors.New(cors.Options{
		AllowedOrigins: engine.options.Cors,
	})

	tokenAuth = jwtauth.New("HS256", []byte(opts.JWTToken), nil)
	ctx := engine.pool.Get().(*Context)

	engine.Use(Recovery)
	engine.Use(Logger)
	engine.Use(crs.Handler)
	engine.Use(tokenAuth.Verifier)
	engine.Use(ctx.Authenticator)

	return engine
}

// Use add middlewares to be used on applicaton
func (c *Engine) Use(middleware Middleware) {
	c.Router.mux.Use(middleware)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (c *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c.Router.mux.ServeHTTP(res, req)
}

// Run run the http server.
func (c *Engine) Run(port int) error {

	addr := fmt.Sprintf(":%d", port)
	l.Printf("Starting server on port [%d]\n", port)
	if err := http.ListenAndServe(addr, c); err != nil {
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
