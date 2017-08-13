package minion

import "github.com/go-chi/chi"

// ByGet shortcut to chi.URLParam
// returns the url parameter from a http.Request object.
func (c *Context) ByGet(name string) string {
	return chi.URLParam(c.Req, name)

}
