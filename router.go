package minion

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi"
)

type Router struct {
	mux       *chi.Mux
	app       *App
	namespace string
}

func (r *Router) handle(method string, relativePath string, handler HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)

	fn := func(rw http.ResponseWriter, req *http.Request) {
		ctx := r.app.createContext(rw, req, handler)
		ctx.Execute()
		r.app.reuseContext(ctx)
	}

	switch method {
	case "ALL":
		r.mux.HandleFunc(namespace, fn)
	case "GET":
		r.mux.Get(namespace, fn)
	case "POST":
		r.mux.Post(namespace, fn)
	case "PUT":
		r.mux.Put(namespace, fn)
	case "PATCH":
		r.mux.Patch(namespace, fn)
	case "DELETE":
		r.mux.Delete(namespace, fn)
	case "OPTIONS":
		r.mux.Options(namespace, fn)
	case "HEAD":
		r.mux.Head(namespace, fn)
	}
}

// Post handle the POST requests
func (r *Router) Post(relativePath string, handler HandlerFunc) {
	r.handle("POST", relativePath, handler)
}

// Get handle the GET requets
func (r *Router) Get(relativePath string, handler HandlerFunc) {
	r.handle("GET", relativePath, handler)
}

// Delete handle the DELETE requests
func (r *Router) Delete(relativePath string, handler HandlerFunc) {
	r.handle("DELETE", relativePath, handler)
}

// Patch handle the PATCH requests
func (r *Router) Patch(relativePath string, handler HandlerFunc) {
	r.handle("PATCH", relativePath, handler)
}

// Put handle the PUT requests
func (r *Router) Put(relativePath string, handler HandlerFunc) {
	r.handle("PUT", relativePath, handler)
}

// Options handle the OPTIONS requests
func (r *Router) Options(relativePath string, handler HandlerFunc) {
	r.handle("OPTIONS", relativePath, handler)
}

// Head handle the HEAD requests
func (r *Router) Head(relativePath string, handler HandlerFunc) {
	r.handle("HEAD", relativePath, handler)
}

// Handle handle all paths to any http method
func (r *Router) Handle(relativePath string, handler HandlerFunc) {
	r.handle("ALL", relativePath, handler)
}

func (r *Router) StaticServer(path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.mux.Get(path, http.RedirectHandler(path+"/", http.StatusTemporaryRedirect).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.mux.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (r *Router) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return r.namespace
	}
	namespace := path.Join(r.namespace, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(namespace) != '/'
	if appendSlash {
		return namespace + "/"
	}
	return namespace
}
