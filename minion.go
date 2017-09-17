package minion

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

var l = log.New(os.Stdout, "[minion] ", 0)

type HandlerFunc func(*Context)
type Middleware func(http.Handler) http.Handler

var tokenAuth *jwtauth.JwtAuth

type App struct {
	*Router
	c       *Context
	pool    sync.Pool
	options Options
}

// Options defines the options to start the API
type Options struct {
	Cors                  []string
	JWTToken              string
	UnauthenticatedRoutes []string
	Namespace             string
	Authenticator         func(next http.Handler) http.Handler
}

func New(opts Options) *App {
	namespace := opts.Namespace
	if len(namespace) == 0 {
		namespace = "/"
	}

	app := &App{}
	app.Router = &Router{
		app: app,
		mux: chi.NewRouter(),
	}

	app.options = opts
	app.pool.New = func() interface{} {
		ctx := &Context{
			app: app,
			render: render.New(render.Options{
				Layout: "layout",
			}),
		}
		return ctx
	}

	return app
}

// Classic returns a new Engine instance with basic middlewares
// Recovery, Logger, CORS and JWT
func Classic(opts Options) *App {
	app := New(opts)

	crs := cors.New(cors.Options{
		AllowedOrigins:   app.options.Cors,
		AllowedHeaders:   []string{"Authorization", "Origin", "X-Requested-With", "Content-Type", "Accept"},
		AllowCredentials: true,
	})

	tokenAuth = jwtauth.New("HS256", []byte(opts.JWTToken), nil)
	ctx := app.pool.Get().(*Context)

	app.Use(middleware.Recoverer)
	app.Use(Logger)
	app.Use(crs.Handler)
	app.Use(jwtauth.Verifier(tokenAuth))
	if opts.Authenticator != nil {
		app.Use(opts.Authenticator)
	} else {
		app.Use(ctx.Authenticator)
	}

	return app
}

func (app *App) Use(md Middleware) {
	app.Router.mux.Use(md)
}

func (app *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	app.Router.mux.ServeHTTP(rw, req)
}

func (app *App) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port [%d]\n", port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(listen, app)
}

func (app *App) reuseContext(ctx *Context) {
	app.pool.Put(ctx)
}
