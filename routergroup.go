package minion

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

// RouterGroup TODO
type RouterGroup struct {
	Handlers     []HandlerFunc
	absolutePath string
	engine       *Engine
}

// Use Adds middlewares to the group
func (c *RouterGroup) Use(middlewares ...HandlerFunc) {
	c.Handlers = append(c.Handlers, middlewares...)
}

// Group Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (c *RouterGroup) Group(relativePath string, fn func(*RouterGroup), handlers ...HandlerFunc) *RouterGroup {
	router := &RouterGroup{
		Handlers:     c.combineHandlers(handlers),
		absolutePath: c.calculateAbsolutePath(relativePath),
		engine:       c.engine,
	}
	fn(router)
	return router
}

// Handle registers a new request handle and middlewares with the given path and method.
// The last handler should be the real handler, the other ones should be middlewares that can and should be shared among different routes.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (c *RouterGroup) Handle(httpMethod, relativePath string, handlers []HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Handle(httpMethod, absolutePath, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := c.engine.createContext(w, req, params, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (c *RouterGroup) Post(relativePath string, handlers ...HandlerFunc) {
	c.Handle("POST", relativePath, handlers)
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (c *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) {
	c.Handle("GET", relativePath, handlers)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (c *RouterGroup) Delete(relativePath string, handlers ...HandlerFunc) {
	c.Handle("DELETE", relativePath, handlers)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (c *RouterGroup) Patch(relativePath string, handlers ...HandlerFunc) {
	c.Handle("PATCH", relativePath, handlers)
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (c *RouterGroup) Put(relativePath string, handlers ...HandlerFunc) {
	c.Handle("PUT", relativePath, handlers)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (c *RouterGroup) Options(relativePath string, handlers ...HandlerFunc) {
	c.Handle("OPTIONS", relativePath, handlers)
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (c *RouterGroup) Head(relativePath string, handlers ...HandlerFunc) {
	c.Handle("HEAD", relativePath, handlers)
}

// Any is a shortcut for all request methods
func (c *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) {
	c.Handle("GET", relativePath, handlers)
	c.Handle("POST", relativePath, handlers)
	c.Handle("PUT", relativePath, handlers)
	c.Handle("PATCH", relativePath, handlers)
	c.Handle("HEAD", relativePath, handlers)
	c.Handle("OPTIONS", relativePath, handlers)
	c.Handle("DELETE", relativePath, handlers)
	c.Handle("CONNECT", relativePath, handlers)
	c.Handle("TRACE", relativePath, handlers)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use : router.Static("/static", "/var/www")
func (c *RouterGroup) Static(path, dir string) {
	if lastChar(path) != '/' {
		path += "/"
	}
	path += "*filepath"
	c.engine.router.ServeFiles(path, http.Dir(dir))
}

func (c *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(c.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, c.Handlers...)
	return append(mergedHandlers, handlers...)
}

func (c *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return c.absolutePath
	}
	absolutePath := path.Join(c.absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(absolutePath) != '/'
	if appendSlash {
		return absolutePath + "/"
	}
	return absolutePath
}
