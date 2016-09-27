package minion

import (
	"math"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/unrolled/render"
)

// Context the context of each request
type Context struct {
	Writer  http.ResponseWriter
	Req     *http.Request
	Session Session
	Engine  *Engine
	render  *render.Render
	// writer   writer
	handlers []HandlerFunc
	index    int8
}

const (
	abortIndex = math.MaxInt8 / 2
)

// Next should be used only in the middlewares.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Abort Forces the system to do not continue calling the pending handlers in the chain.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Redirect returns a HTTP redirect to the specific location. default for 302
func (c *Context) Redirect(status int, location string) {
	c.SetHeader("Location", location)
	if status > 0 {
		http.Redirect(c.Writer, c.Req, location, status)
	} else {
		http.Redirect(c.Writer, c.Req, location, http.StatusTemporaryRedirect)
	}
}

// JSON Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(status int, data interface{}) {
	c.render.JSON(c.Writer, status, data)
}

// JSONP Serializes the given struct as JSONP into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/javascript".
func (c *Context) JSONP(status int, data interface{}, callback string) {
	c.render.JSONP(c.Writer, status, callback, data)
}

// Text Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) Text(status int, data string) {
	c.render.Text(c.Writer, status, data)
}

// HTML renders the html template and sets the Content-Type to "text/html".
func (c *Context) HTML(status int, tmpl string, data interface{}) {
	c.render.HTML(c.Writer, status, tmpl, data)
}

// MarshalOnePayload marshal the struct and return as jsonaapi
func (c *Context) MarshalOnePayload(status int, model interface{}) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Set("Content-Type", "application/vnd.api+json")
	if err := jsonapi.MarshalOnePayload(c.Writer, model); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// MarshalManyPayload marshal the struct and return as jsonaapi
func (c *Context) MarshalManyPayload(status int, models []interface{}) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Set("Content-Type", "application/vnd.api+json")
	if err := jsonapi.MarshalManyPayload(c.Writer, models); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Engine) createContext(rw http.ResponseWriter, req *http.Request, handlers []HandlerFunc) *Context {
	ctx := c.pool.Get().(*Context)
	ctx.Writer = rw
	ctx.Req = req
	ctx.handlers = handlers
	ctx.index = -1
	return ctx
}

func (c *Engine) reuseContext(ctx *Context) {
	c.pool.Put(ctx)
}
