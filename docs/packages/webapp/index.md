# webapp

Welcome to the webapp user guide. Here we will guide you through some
practical examples of how to interact with the API.

A webapp.Application is what most people will end up using, this struct
is just container of other components exposed by the "go-toolkit",
constructed with sane defaults, compliant with the behavior that's expected
from an application working in the web environment.

## Install

    go get -u github.com/pomelo-la/go-toolkit/webapp

## Create a Application

```go
package main

import (
	"context"
	"log"
	
	"github.com/pomelo-la/go-toolkit/webapp"
)

func main() {
	ctx := context.Background()
	app, err := webapp.New("go-toolkit-webapp")
	if err != nil {
		log.Fatalf("error initializing web app")
	}

	// run webapp
	// under the hood shutdown app is to gracefully
	if err := app.Run(); err != nil {
		app.Logger.Error(ctx, "error during run app")
	}
}
```

Add endpoints to your app

```go hl_lines="6 19-21"
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/webapp"
)

func main() {
	ctx := context.Background()
	app, err := webapp.New("go-toolkit-webapp")
	if err != nil {
		log.Fatalf("error initializing web app")
	}
	
	app.Router.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		return httprouter.RespondJSON(w, http.StatusOK, nil)
	})

	// run webapp
	// under the hood shutdown app is to gracefully
	if err := app.Run(); err != nil {
		app.Logger.Error(ctx, "error during run app")
	}
}
```

## Configuration new app

### Add global middlewares

```go
package middlewares

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Content-Type, Content-Length, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, HEAD, DELETE, PATCH")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

!!! note

    Is just an example, please adapt the cors middleware to your needs.

```go hl_lines="11 18"
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pomelo-la/go-toolkit/httprouter"
	"github.com/pomelo-la/go-toolkit/webapp"

	"github.com/my-org/my-awesome-app/middlewares"
)

func main() {
	ctx := context.Background()
	
	app, err := webapp.New("go-toolkit-webapp", 
		webapp.WithGlobalMiddlewares(Cors))
	if err != nil {
		log.Fatalf("error initializing web app")
	}
	
	app.Router.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		return httprouter.RespondJSON(w, http.StatusOK, nil)
	})

	// run webapp
	// under the hood shutdown app is to gracefully
	if err := app.Run(); err != nil {
		app.Logger.Error(ctx, "error during run app")
	}
}
```
