# httprouter

Welcome to the httprouter user guide. Here we will guide you through 
some practical examples of how to interact with the API.

httprouter is built around go-chi, we pursue its philosophy and 
consider `func(w http.ResponseWriter, r *http.Request)` to be 
more extensible and conform to standard http handlers.

!!! info

    One of the most important changes we can find is the handler type.

    `type Handler func(w http.ResponseWriter, r *http.Request) error`

We believe that this new handler signature allows us to handle more 
subtly what a handler returns and to wrap and return errors 
more elegantly to the client.

## Install

    go get -u github.com/pomelo-la/go-sdk/httprouter

## Running a simple server

```go
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/pomelo-la/go-toolkit/httprouter"
)

type Greet struct {
	PkgName string `json:"pkg_name"`
}

func main() {
	r := httprouter.New()
	
	r.Get("/hello/{go-toolkit}", func(w http.ResponseWriter, r *http.Request) error {
		urlParams := httprouter.URLParam(r, "go-toolkit")

		return httprouter.RespondJSON(w, http.StatusOK, Greet{PkgName: urlParams})
	})
	
	ln, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server at %s", ln.Addr().String())

	if err := httprouter.Run(ln, httprouter.DefaultTimeouts, w); err != nil {
		log.Fatal(err)
	}
}
```
Now run your simple API

    go run main.go

Execute curl to call the API

	curl http://localhost:9090/hello/httprouter
	Expected Output:{"pkg_name":"httprouter"}

### Let us explain step by step what we did.
First we instantiate the router that will be used for registering 
HTTP handlers.

    r := httprouter.New()

Then we register a GET handler. This handler receives a string by 
URI param and returns a JSON with a greeting message.

    r.Get("/hello/{go-toolkit}", func(w http.ResponseWriter, r *http.Request) error {
		urlParam := httprouter.URLParam(r, "go-toolkit")

		return httprouter.RespondJSON(w, http.StatusOK, Greet{PkgName: urlParam})
	})

Then we create the listener that will be pass in to the underlying 
http.Server for attending incoming requests.

    ln, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

Finally, we Run blocks and listen for an interrupt or terminate signal from the OS.
After signalling, a graceful shutdown is triggered without affecting any 
live connections/clients connected to the server. 
It will complete executing all the active/live requests before shutting down.

    if err := httprouter.Run(ln, httprouter.DefaultTimeouts, w); err != nil {
		log.Fatal(err)
	}
