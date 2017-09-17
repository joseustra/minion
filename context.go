package minion

import (
	"net/http"

	"github.com/unrolled/render"
)

type Context struct {
	writer  http.ResponseWriter
	Req     *http.Request
	render  *render.Render
	handler HandlerFunc
	app     *App
}

func (c *Context) Execute() {
	c.handler(c)
}

func (app *App) createContext(rw http.ResponseWriter, req *http.Request, handler HandlerFunc) *Context {
	ctx := app.pool.Get().(*Context)
	ctx.writer = rw
	ctx.Req = req
	ctx.handler = handler
	return ctx
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) JSON(status int, v interface{}) {
	c.render.JSON(c.writer, status, v)
}

func (c *Context) Text(status int, v string) {
	c.render.Text(c.writer, status, v)
}

func (c *Context) HTML(status int, tmpl string, v interface{}) {
	c.render.HTML(c.writer, status, tmpl, v)
}

// Redirect returns a HTTP redirect to the specific location. default for 302
func (c *Context) Redirect(status int, location string) {
	c.SetHeader("Location", location)
	if status > 0 {
		http.Redirect(c.writer, c.Req, location, status)
	} else {
		http.Redirect(c.writer, c.Req, location, http.StatusTemporaryRedirect)
	}
}
