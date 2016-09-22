package minion

import (
	"net/http"
	"path"
)

// Router TODO
type Router struct {
	Handlers  []HandlerFunc
	namespace string
	engine    *Engine
}

// Use Adds middlewares to the group
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.Handlers = append(c.Handlers, middlewares...)
}

func (c *Router) handleContext(w http.ResponseWriter, req *http.Request, handlers []HandlerFunc) {
	ctx := c.engine.createContext(w, req, handlers)
	ctx.Next()
	ctx.Writer.WriteHeaderNow()
	c.engine.reuseContext(ctx)
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (c *Router) Post(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Post(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (c *Router) Get(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Get(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (c *Router) Delete(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Delete(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (c *Router) Patch(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Patch(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (c *Router) Put(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Put(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (c *Router) Options(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Options(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (c *Router) Head(relativePath string, handlers ...HandlerFunc) {
	namespace := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Head(namespace, func(w http.ResponseWriter, req *http.Request) {
		c.handleContext(w, req, handlers)
	})
}

// Static serves files from the given file system root.
// use : router.Static("/static", "/var/www")
func (c *Router) Static(path, dir string) {
	c.engine.router.FileServer(path, http.Dir(dir))
}

func (c *Router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(c.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, c.Handlers...)
	return append(mergedHandlers, handlers...)
}

func (c *Router) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return c.namespace
	}
	namespace := path.Join(c.namespace, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(namespace) != '/'
	if appendSlash {
		return namespace + "/"
	}
	return namespace
}
