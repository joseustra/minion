package minion_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goware/jwtauth"
	"github.com/stretchr/testify/assert"
	"github.com/ustrajunior/minion"
	"github.com/ustrajunior/minion/tst"
)

func TestGet(t *testing.T) {
	john := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		"John Doe",
		"john@doe.com",
	}

	foo := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		"Foo Bar",
		"foo@bar.com",
	}

	usersHandler := func(ctx *minion.Context) {
		users := struct {
			Users []interface{} `json:"users"`
		}{
			[]interface{}{john, foo},
		}
		ctx.JSON(200, users)
	}

	userHandler := func(ctx *minion.Context) {
		id := ctx.ByGet("id")
		doe := struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			id,
			"John Doe",
			"john@doe.com",
		}
		ctx.JSON(200, doe)
	}

	m := minion.New(minion.Options{})

	m.Get("/users", usersHandler)
	m.Get("/user/{id}", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	status, body = tst.Request(t, ts, "GET", "/users", nil, nil)
	assert.Equal(t, 200, status)

	j = `{"users":[{"name":"John Doe","email":"john@doe.com"},{"name":"Foo Bar","email":"foo@bar.com"}]}`
	assert.Equal(t, j, body)

	status, body = tst.Request(t, ts, "GET", "/user/1", nil, nil)
	assert.Equal(t, 200, status)

	j = `{"id":"1","name":"John Doe","email":"john@doe.com"}`
	assert.Equal(t, j, body)
}

func TestPost(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	usersHandler := func(ctx *minion.Context) {
		user := &struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{}

		ctx.GetResource(user)

		user.ID = "1"

		ctx.JSON(200, user)
	}

	m.Post("/users", usersHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	payload := `{"name":"John","email":"john@doe.com"}`
	status, body = tst.Request(t, ts, "POST", "/users", nil, bytes.NewBuffer([]byte(payload)))
	assert.Equal(t, 200, status)

	j = `{"id":"1","name":"John","email":"john@doe.com"}`
	assert.Equal(t, j, body)
}

func TestPatch(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	type user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	userHandler := func(ctx *minion.Context) {
		dbUser := &user{"1", "John", "foo@bar.com"}
		user := new(user)
		ctx.GetResource(user)

		dbUser.Email = user.Email
		ctx.JSON(200, dbUser)
	}

	m.Patch("/users/{id}", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	payload := `{"email":"john@doe.com"}`
	status, body = tst.Request(t, ts, "PATCH", "/users/1", nil, bytes.NewBuffer([]byte(payload)))
	assert.Equal(t, 200, status)

	j = `{"id":"1","name":"John","email":"john@doe.com"}`
	assert.Equal(t, j, body)
}

func TestPut(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	type user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	userHandler := func(ctx *minion.Context) {
		dbUser := &user{"1", "John", "foo@bar.com"}
		user := new(user)
		ctx.GetResource(user)

		dbUser.Name = user.Name
		dbUser.Email = user.Email
		ctx.JSON(200, dbUser)
	}

	m.Put("/users/{id}", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	payload := `{"name":"John Doe","email":"john@doe.com"}`
	status, body = tst.Request(t, ts, "PUT", "/users/1", nil, bytes.NewBuffer([]byte(payload)))
	assert.Equal(t, 200, status)

	j = `{"id":"1","name":"John Doe","email":"john@doe.com"}`
	assert.Equal(t, j, body)
}

func TestOptions(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	userHandler := func(ctx *minion.Context) {
		ctx.Text(200, "")
	}

	m.Options("/users", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var body string
	var status int

	status, body = tst.Request(t, ts, "OPTIONS", "/users", nil, nil)
	assert.Equal(t, 200, status)

	assert.Equal(t, "", body)
}

func TestHead(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	userHandler := func(ctx *minion.Context) {
		ctx.Text(200, "")
	}

	m.Head("/users", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var body string
	var status int

	status, body = tst.Request(t, ts, "HEAD", "/users", nil, nil)
	assert.Equal(t, 200, status)

	assert.Equal(t, "", body)
}

func TestDelete(t *testing.T) {
	m := minion.New(minion.Options{
		UnauthenticatedRoutes: []string{"*"}},
	)

	userHandler := func(ctx *minion.Context) {
		j := struct {
			Message string `json:"message"`
		}{
			"ok",
		}
		ctx.JSON(200, j)
	}

	m.Delete("/users/{id}", userHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	status, body = tst.Request(t, ts, "DELETE", "/users/1", nil, nil)
	assert.Equal(t, 200, status)

	j = `{"message":"ok"}`
	assert.Equal(t, j, body)
}

func TestJWTAuthentication(t *testing.T) {
	m := minion.Classic(minion.Options{JWTToken: "123"})

	usersHandler := func(ctx *minion.Context) {
		j := struct {
			Message string `json:"message"`
		}{
			"ok",
		}
		ctx.JSON(200, j)
	}

	m.Get("/users", usersHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	tokenString, _ := minion.CreateJWTToken(nil)

	h := http.Header{}
	h.Set("Authorization", "BEARER "+tokenString)

	status, body = tst.Request(t, ts, "GET", "/users", h, nil)
	assert.Equal(t, 200, status)

	j = `{"message":"ok"}`
	assert.Equal(t, j, body)
}

func TestJWTAuthenticationWrongToken(t *testing.T) {
	m := minion.Classic(minion.Options{JWTToken: "123"})

	usersHandler := func(ctx *minion.Context) {
		j := struct {
			Message string `json:"message"`
		}{
			"ok",
		}
		ctx.JSON(200, j)
	}

	m.Get("/users", usersHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	tokenAuth := jwtauth.New("HS256", []byte("wrong"), nil)
	_, tokenString, _ := tokenAuth.Encode(nil)

	h := http.Header{}
	h.Set("Authorization", "BEARER "+tokenString)

	status, body = tst.Request(t, ts, "GET", "/users", h, nil)
	assert.Equal(t, 401, status)

	j = `{"status":401,"message":"Unauthorized"}`
	assert.Equal(t, j, body)
}

func TestJWTAuthenticationUnauthenticatedPath(t *testing.T) {
	m := minion.Classic(minion.Options{
		JWTToken:              "123",
		UnauthenticatedRoutes: []string{"/unauthenticated"},
	})

	usersHandler := func(ctx *minion.Context) {
		j := struct {
			Message string `json:"message"`
		}{
			"ok",
		}
		ctx.JSON(200, j)
	}

	m.Get("/unauthenticated", usersHandler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	var j, body string
	var status int

	status, body = tst.Request(t, ts, "GET", "/unauthenticated", nil, nil)
	assert.Equal(t, 200, status)

	j = `{"message":"ok"}`
	assert.Equal(t, j, body)
}
