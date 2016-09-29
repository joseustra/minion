package minion

import (
	"net/http"
	"path"

	"github.com/pressly/chi"
)

// Router TODO
type Router struct {
	Handlers  []HandlerFunc
	namespace string
	engine    *Engine
	mux       *chi.Mux
}

func (r *Router) handle(method string, relativePath string, handlers HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	fn := func(rw http.ResponseWriter, req *http.Request) {
		ctx := r.engine.createContext(rw, req, []HandlerFunc{handlers})
		ctx.Next()
		r.engine.reuseContext(ctx)
	}

	switch method {
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
func (r *Router) Post(relativePath string, handlers HandlerFunc) {
	r.handle("POST", relativePath, handlers)
}

// Get handle the GET requets
func (r *Router) Get(relativePath string, handlers HandlerFunc) {
	r.handle("GET", relativePath, handlers)
}

// Delete handle the DELETE requests
func (r *Router) Delete(relativePath string, handlers HandlerFunc) {
	r.handle("DELETE", relativePath, handlers)
}

// Patch handle the PATCH requests
func (r *Router) Patch(relativePath string, handlers HandlerFunc) {
	r.handle("PATCH", relativePath, handlers)
}

// Put handle the PUT requests
func (r *Router) Put(relativePath string, handlers HandlerFunc) {
	r.handle("PUT", relativePath, handlers)
}

// Options handle the OPTIONS requests
func (r *Router) Options(relativePath string, handlers HandlerFunc) {
	r.handle("OPTIONS", relativePath, handlers)
}

// Head handle the HEAD requests
func (r *Router) Head(relativePath string, handlers HandlerFunc) {
	r.handle("HEAD", relativePath, handlers)
}

// Static serves files from the given file system root.
func (r *Router) Static(path, dir string) {
	r.mux.FileServer(path, http.Dir(dir))
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
