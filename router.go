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

// Post is a shortcut for router.Handle("POST", path, handle)
func (r *Router) Post(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Post(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (r *Router) Get(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Get(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) Delete(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Delete(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) Patch(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Patch(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) Put(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Put(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) Options(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Options(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) Head(relativePath string, handlers ...HandlerFunc) {
	namespace := r.calculateAbsolutePath(relativePath)
	r.mux.Head(namespace, func(w http.ResponseWriter, req *http.Request) {
		r.handleContext(w, req, handlers)
	})
}

// Static serves files from the given file system root.
// use : router.Static("/static", "/var/www")
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
