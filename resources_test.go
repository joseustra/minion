package minion_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ustrajunior/minion"
	"github.com/ustrajunior/minion/tst"
)

func TestGetResource(t *testing.T) {
	user := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		"1",
		"John Doe",
		"john@doe.com",
	}

	body := bytes.NewBufferString(`{"id":"1","name":"John Doe","email":"john@doe.com"}`)

	m := minion.New(minion.Options{})
	m.Post("/users", nil)

	ts := httptest.NewServer(m)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/users", body)

	var resource struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	ctx := &minion.Context{Req: req}
	err := ctx.GetResource(&resource)

	tst.AssertNoError(t, err)
	tst.AssertEqual(t, user, resource)
}
