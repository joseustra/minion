package minion

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type routerParams struct {
	req    *http.Request
	params httprouter.Params
}

func (c *routerParams) ByGet(name string) string {
	val := c.params.ByName(name)
	if val == "" {
		val = c.req.URL.Query().Get(name)
	}
	return val
}

func (c *routerParams) ByPost(name string) string {
	return c.req.FormValue(name)
}

func (c *routerParams) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.req.FormFile(name)
}

func (c *routerParams) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.req.Body)
	err := decoder.Decode(obj)
	return err
}

func (c *routerParams) JSON() *jsonParams {
	defer c.req.Body.Close()

	data, _ := ioutil.ReadAll(c.req.Body)
	objJSON := &jsonParams{data: map[string]interface{}{}}
	objJSON.source = string(data)
	json.Unmarshal(data, &objJSON.data)

	return objJSON
}

type jsonParams struct {
	source string
	data   map[string]interface{}
}

func (c *jsonParams) Get(name string) interface{} {
	if len(c.data) == 0 || c.data[name] == nil {
		return ""
	}
	return c.data[name]
}

func (c *jsonParams) GetString(name string) string {
	if len(c.data) == 0 || c.data[name] == nil {
		return ""
	}
	return toString(c.data[name])
}

func (c *jsonParams) GetInt32(name string) int32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toInt32(c.data[name])
}

func (c *jsonParams) GetUInt32(name string) uint32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toUint32(c.data[name])
}

func (c *jsonParams) GetFloat32(name string) float32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toFloat32(c.data[name])
}

func (c *jsonParams) GetFloat64(name string) float64 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toFloat64(c.data[name])
}

func (c *jsonParams) String() string {
	return c.source
}
