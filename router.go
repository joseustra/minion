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

func (r *Router) handleContext(w http.ResponseWriter, req *http.Request, handlers []HandlerFunc) {
	ctx := r.engine.createContext(w, req, handlers)
	ctx.Next()
	r.engine.reuseContext(ctx)
}

// Post handle the POST requests
func (r *Router) Post(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Post(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Get handle the GET requets
func (r *Router) Get(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Get(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Delete handle the DELETE requests
func (r *Router) Delete(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Delete(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Patch handle the PATCH requests
func (r *Router) Patch(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Patch(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Put handle the PUT requests
func (r *Router) Put(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Put(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Options handle the OPTIONS requests
func (r *Router) Options(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Options(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Head handle the HEAD requests
func (r *Router) Head(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Head(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
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
