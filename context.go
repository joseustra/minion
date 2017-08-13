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

func (c *Context) JSON(status int, v interface{}) {
	c.render.JSON(c.writer, status, v)
}

func (c *Context) Text(status int, v string) {
	c.render.Text(c.writer, status, v)
}
