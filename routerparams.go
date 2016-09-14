package minion

import (
	"mime/multipart"

	"github.com/pressly/chi"
)

func (c *Context) ByGet(name string) string {
	return chi.URLParam(c.Req, name)
}

func (c *Context) ByQuery(name string) string {
	values := c.Req.URL.Query()
	if len(values) == 0 {
		return ""
	}

	return values[name][0]
}

func (c *Context) ByPost(name string) string {
	return c.Req.FormValue(name)
}

func (c *Context) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.Req.FormFile(name)
}
