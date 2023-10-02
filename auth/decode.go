package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrRequestNotAcceptable indicates that the HTTP request is not acceptable.
	ErrRequestNotAcceptable = errors.New("req not acceptable")
	// ErrUnauthorized indicates that the HTTP request is unauthorized.
	ErrUnauthorized = errors.New("req unauthorized")
	// ErrTypeAssertionsClaims indicates that the type assertions had an error.
	ErrTypeAssertionsClaims = errors.New("invalid claims type assertions")
	// ErrMissingAreas indicates that the token are is empty.
	ErrMissingAreas = errors.New("the areas in the claims token must not be empty")
	// ErrMissingEmail indicates that the token email is empty.
	ErrMissingEmail = errors.New("the email in the claims token must not be empty")
)

const (
	// Role specifies the role of the user or service
	Role = "role"
	// Owner specifies the owner of the token.
	Owner = "owner"
	// BusinessUnits specifies the business units of Pomelo company
	BusinessUnits = "business-units"
	// ServiceContextAPI specifies the service
	ServiceContextAPI = "servicecontextapi"
)

// Claims the JWT Claims Set represents a JSON object whose members are the claims conveyed by the JWT.
//
// RFC for more info https://datatracker.ietf.org/doc/html/rfc7519#section-4
type Claims struct {
	// Area specifies the business units of Pomelo company
	Area []string `json:"area"`
	// Email specifies the email of the user or business unit
	Email string `json:"email"`
	// Role specifies the role of the user or service
	Role []string `json:"role"`
	// ServiceName specifies the name of the service
	ServiceName string `json:"servicename"`
	jwt.RegisteredClaims
}

// DecodeToken provides a logic to decode Pomelo token.
func DecodeToken(r *http.Request) (*Claims, error) {
	tokenHeader := r.Header.Get("X-Auth-Token")
	if tokenHeader == "" {
		return nil, ErrRequestNotAcceptable
	}

	token, err := ensureValidToken(tokenHeader)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTypeAssertionsClaims
	}
	if len(claims.Area) == 0 {
		return nil, ErrMissingAreas
	}
	if claims.Email == "" && !claims.VerifyAudience(ServiceContextAPI, true) {
		return nil, ErrMissingEmail
	}

	r.Header.Set(BusinessUnits, strings.Join(claims.Area, ","))
	r.Header.Set(Owner, claims.Email)

	if len(claims.Role) > 0 {
		r.Header.Set(Role, claims.Role[0])
	}
	if claims.VerifyAudience(ServiceContextAPI, true) {
		serviceName := claims.ServiceName
		r.Header.Set(Role, "service_"+serviceName)
		r.Header.Set(Owner, serviceName)
	}

	return claims, nil
}

func ensureValidToken(tokenHeader string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenHeader, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("CONTEXT_API_PUBLIC_KEY")))
	})
	switch {
	case errors.Is(err, jwt.ErrSignatureInvalid):
		return nil, ErrRequestNotAcceptable
	case err != nil:
		return nil, ErrUnauthorized
	}

	if !token.Valid {
		return nil, ErrUnauthorized
	}

	return token, nil
}
