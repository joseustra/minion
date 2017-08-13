# Minion

## A micro RESTfull "framework" for Go apps.

**Minion** is a very small library to help you start a new Go App. It will handle the basics so you can focus on developing your application.

Out of the box, it includes a Router, a template engine, JWT, logger, a recovery middleware, etc.

## Status

The source code was rewritten to be simpler and many functions are still in development. Do **not** use in production.

## How to install

```go
go get -u github.com/ustrajunior/minion
```

## Usage

A very simple app.

```go
package main

import (
	"github.com/ustrajunior/minion"
)

func main() {
	m := minion.New(minion.Options{})

	m.Get("/", homeHandler)
	m.Run(8800)
}

func homeHandler(ctx *minion.Context) {
	project := struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		"Minion",
		"https://github.com/ustrajunior/minion",
	}

	ctx.JSON(http.StatusOK, project)
}
```

## External packages used to create Minion

* [go-chi/chi](https://github.com/go-chi/chi)
* [unrolled/render](https://github.com/unrolled/render)

## License
Minion is licensed under the MIT license.
