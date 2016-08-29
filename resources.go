package minion

import (
	"encoding/json"
	"io/ioutil"
)

// GetResource unmarshal the request body and return the resource
func (c *Context) GetResource(resource interface{}) error {
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resource)
}
