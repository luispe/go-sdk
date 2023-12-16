# Middleware

!!! note
    From go-chi definition

    Middleware performs some specific function on the HTTP 
    request or response at a specific stage in the HTTP 
    pipeline before or after the user defined controller. 
    Middleware is a design pattern to eloquently add 
    cross-cutting concerns like logging, handling authentication 
    without having many code contact points.

## Declare middlewares

Here is an example of a standard net/http middleware
where we assign a context key "user" the value of "123". 
This middleware sets a hypothetical user identifier on 
the request context and calls the next handler in the chain.

```go
// HTTP middleware setting a value on the request context
func MyMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // create new context from `r` request context, and assign key `"user"`
    // to value of `"123"`
    ctx := context.WithValue(r.Context(), "user", "123")

    // call the next handler in the chain, passing the response writer and
    // the updated request object with the new context value.
    //
    // note: context.Context values are nested, so any previously set
    // values will be accessible as well, and the new `"user"` key
    // will be accessible from this point forward.
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}
```

We can now take these values from the context in our Handlers like this:

```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    // here we read from the request context and fetch out `"user"` key set in
    // the MyMiddleware example above.
    user := r.Context().Value("user").(string)

    // respond to the client
    w.Write([]byte(fmt.Sprintf("hi %s", user)))
}
```

The following example is a middleware of the webapp package, 
it is shared in order to have an example where middleware 
receives parameters

```go
// log decorates the request context with the given logger, accessible via
// the go-toolkit log methods with context.
func logMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(r.Context(), "request started",
				slog.String("method", r.Method),
				slog.String("path", path),
				slog.String("remoteaddr", r.RemoteAddr),
			)

			next.ServeHTTP(w, r)
		})
	}
}
```

## Working with middlewares

httprouter offers, like go-chi, two ways to inject middleware, 
global or inline, `r.Use()` and `r.With()` respectively.

By go-chi design a global middleware can only be declared before 
declaring the endpoints, otherwise you will get a router panic, 
here we will explore some scenarios for working with middlewares.

!!! danger

    The following example will return a panic

```go
r := httprouter.NewRouter()

r.Get("/articles/{year}-{month}", getArticle)

r.Use(MyMiddleware)

func getArticle(w http.ResponseWriter, r *http.Request) error {
    yearParam := httprouter.URLParam(r, "year")
    monthParam := httprouter.URLParam(r, "month")
    // Add your logic
})
```

As we can see we first declare the endpoint and then inject 
the middleware, let's correct it.

!!! success

    The following example will work

```go
r := httprouter.NewRouter()

r.Use(MyMiddleware)

r.Get("/articles/{year}-{month}", getArticle)

func getArticle(w http.ResponseWriter, r *http.Request) error {
    yearParam := httprouter.URLParam(r, "year")
    monthParam := httprouter.URLParam(r, "month")
    // Add your logic
})
```

Let's explore the inline middleware scenario.

```go hl_lines="6"
r := httprouter.NewRouter()

r.Use(MyMiddleware)

r.Get("/articles/{year}-{month}", getArticle)
r.With(authDelete).Delete("/articles/{id}", deleteArticle)

func getArticle(w http.ResponseWriter, r *http.Request) error {
    yearParam := httprouter.URLParam(r, "year")
    monthParam := httprouter.URLParam(r, "month")
    // Add your logic
})

func deleteArticle(w http.ResponseWriter, r *http.Request) error {
    articelID := httprouter.URLParam(r, "id")
    // Add your logic
})
```

In line 6 we can notice that we use middleware only for 
the endpoint of DeleteArticle using

    r.With(authDelete).Delete("/articles/{id}", deleteArticle)


!!! question

    Why does the following scenario work?

```go hl_lines="15"
func main(){
    r := httprouter.NewRouter()
    
    // Public Routes
    r.Group(func(r httprouter.Router) {
        r.Get("/", HelloWorld)
        r.Get("/{AssetUrl}", GetAsset)
        r.Get("/manage/url/{path}", FetchAssetDetailsByURL)
        r.Get("/manage/id/{path}", FetchAssetDetailsByID)
    })

    // Private Routes
    // Require Authentication
    r.Group(func(r httprouter.Router) {
        r.Use(AuthMiddleware)
        r.Post("/manage", CreateAsset)
    })
}
```

As we have mentioned `r.Route() and r.Group()` create 
a new routing instance so a middleware can be declared as we do in

    // Private Routes
    // Require Authentication
    r.Group(func(r httprouter.Router) {
        r.Use(AuthMiddleware)
        r.Post("/manage", CreateAsset)
    })
