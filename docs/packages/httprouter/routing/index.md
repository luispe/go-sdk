# Routing

## Handling HTTP Request Methods

These methods are defined on the `httprouter.Router` as:

    // HTTP-method routing along `pattern`
    Connect(pattern string, h Handler)
    Delete(pattern string, h Handler)
    Get(pattern string, h Handler)
    Head(pattern string, h Handler)
    Options(pattern string, h Handler)
    Patch(pattern string, h Handler)
    Post(pattern string, h Handler)
    Put(pattern string, h Handler)
    Trace(pattern string, h Handler)

Where handler as discussed in the introduction is: 

    type Handler func(w http.ResponseWriter, r *http.Request) error

From here it is basically the same API as [go-chi](https://github.com/go-chi/chi)

## Routing patterns & url parameters

Each routing method accepts a URL pattern and chain of handlers.

The URL pattern supports named params (ie. /users/{userID}) 
and wildcards (ie. /admin/*).

URL parameters can be fetched at runtime by calling
httprouter.URLParam(r, "userID") for named parameters 
and httprouter.URLParam(r, "*") for a wildcard parameter.

### url parameters

```go
r := httprouter.NewRouter()

r.Get("/articles/{year}-{month}", getArticle)

func getArticle(w http.ResponseWriter, r *http.Request) error {
    yearParam := httprouter.URLParam(r, "year")
    monthParam := httprouter.URLParam(r, "month")
    // Add your logic
})
```

as you can see above, the url parameters are defined using the curly 
brackets `{}` with the parameter name in between, as `{year} and {month}`.

When a HTTP request is sent to the server and handled by the router, 
if the URL path matches the format of /articles/{year}-{month}, 
then the `getArticle` function will be called to 
send a response to the client.

### Sub Routers

```go
r := httprouter.NewRouter()

r.Get("/", func(w http.ResponseWriter, r *http.Request) error {
    return httprouter.RespondJSON(w, http.StatusOK, "Hello world!")
})

// Creating a New Router
apiRouter := httprouter.NewRouter()
apiRouter.Get("/articles/{year}-{month}", getArticle)

// Mounting the new Sub Router on the main router
r.Mount("/api", apiRouter)
```

### Another Way of Implementing Sub Routers

```go
r := httprouter.NewRouter()

articleRoutes := r.Route("/articles", 
    func(r httprouter.Router) {})
articleRoutes.Get("/", listArticles)                            // GET /articles
articleRoutes.Get("/{month}-{day}-{year}", listArticlesByDate)  // GET /articles/01-16-2017

articleRoutes.Post("/", createArticle)                          // POST /articles
articleRoutes.Get("/search", searchArticles)                    // GET /articles/search

// Regexp url parameters:
articleRoutes.Get("/{articleSlug:[a-z-]+}", getArticleBySlug)   // GET /articles/home-is-toronto

// Subrouters:
articleRoutesID := articleRoutes.Route("/{articleID}", 
    func(r httprouter.Router) {})
articleRoutesID.Get("/", getArticle)                            // GET /articles/123
articleRoutesID.Put("/", updateArticle)                         // PUT /articles/123
articleRoutesID.Delete("/", deleteArticle)                      // DELETE /articles/123
```

### Routing Groups

You can create Groups in Routers to segregate routes 
using a middleware and some not using a middleware

for example:

```go
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

!!! note
    There is a section dedicated to middleware, 
    soon we will go deeper into middleware.
