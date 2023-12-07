# Other functionalities provided by httprouter

## Bind request

```go
// UserDTO.
type UserDTO struct {
	Name string `validate:"required"`
	Age string `validate:"required"`
}

func (ctrl *RestController) Save(w http.ResponseWriter, r *http.Request) error {
    var userPayload UserDTO
    if err := httprouter.Bind(r, &userPayload); err != nil {
        ctrl.log.Warn(r.Context(), "unprocessable user entity")
        return httprouter.NewError(http.StatusUnprocessableEntity, err.Error())
    }
	// ... another logic
}
```

As we can see we use `httprouter.Bind(r, &userPayload)` 
to validate the payload where `r is the request` and 
`&userPayload` is of type any.

## Response

The RespondJSON method converts a Go value to JSON and sends it to the client.

e.g:
```go
r.Get("/", func(w http.ResponseWriter, r *http.Request) error {
    return httprouter.RespondJSON(w, http.StatusOK, "Hello world!")
})
```

The signature of httprouter.RespondJSON is a:

    func RespondJSON(w http.ResponseWriter, code int, value any) error

## Errors

httprouter offers some methods that will help you to handle 
errors in a unified way without a cognitive load that 
derives from your domain business.

    func NewError(statusCode int, message string) error

or

    func NewErrorf(statusCode int, format string, args ...any) error

Which one to use? It's up to you and what best suits your current problems, 
here are a few basic examples

```go
r.Get("/users/{id}", getUser)

func getUser(w http.ResponseWriter, r *http.Request) error {
    userID := httprouter.URLParam(r, "id")
	if len(userID) < 11 {
		return httprouter.NewError(http.StatusBadRequest, "user id can not be less than 11 characters")
    }
})
```

Or using `NewErrorf`

```go
r.Get("/users/{id}", getUser)

func getUser(w http.ResponseWriter, r *http.Request) error {
    userID := httprouter.URLParam(r, "id")
	if len(userID) < 11 {
		return httprouter.NewErrorf(http.StatusBadRequest, "the following user id: %s is not allowed, it is less than 11 characters long", userID)
    }
})
```
