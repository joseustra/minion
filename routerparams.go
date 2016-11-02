package minion

import (
	"mime/multipart"

	"github.com/pressly/chi"
)

// ByGet shortcut to chi.URLParam
// returns the url parameter from a http.Request object.
func (c *Context) ByGet(name string) string {
	return chi.URLParam(c.Req, name)
}

// ByQuery shortcut to (u *URL) Query()
// parses RawQuery and returns the corresponding values
func (c *Context) ByQuery(name string) string {
	values := c.Req.URL.Query()
	if len(values) == 0 {
		return ""
	}

	return values[name][0]
}

// ByPost shortcut to (r *Request) FormValue
// returns the first value for the named component of the query.
func (c *Context) ByPost(name string) string {
	return c.Req.FormValue(name)
}

// File shortcut to (r *Request) FormFile
// returns the first file for the provided form key
func (c *Context) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.Req.FormFile(name)
}
