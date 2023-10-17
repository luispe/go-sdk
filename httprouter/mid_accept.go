package httprouter

import (
	"mime"
	"net/http"
	"regexp"
	"strings"
)

const allMediaTypes = "*/*"

// AcceptJSON returns middleware that allows requests with "application/json" or an empty Accept header.
// It responds with a NotAcceptable status for other media types.
func AcceptJSON() Middleware {
	return Accept("^application/json$")
}

// Accept creates middleware that allows requests with specified media types.
func Accept(mediaTypes ...string) Middleware {
	compiled := compileMediaTypes(mediaTypes)

	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !isAcceptableRequest(r.Header.Get("Accept"), compiled) {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			handler(w, r)
		}
	}
}

func compileMediaTypes(mediaTypes []string) []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, len(mediaTypes))
	for i, mediaType := range mediaTypes {
		compiled[i] = regexp.MustCompile(mediaType)
	}
	return compiled
}

//revive:disable:cognitive-complexity High complexity score but easy to understand
func isAcceptableRequest(acceptHeader string, mediaTypes []*regexp.Regexp) bool {
	if acceptHeader == "" || acceptHeader == allMediaTypes {
		return true
	}

	acceptedMediaTypes := parseAcceptHeader(acceptHeader)
	if containsWildcard(acceptedMediaTypes...) {
		return true
	}

	for _, accepted := range acceptedMediaTypes {
		for _, t := range mediaTypes {
			if t.MatchString(accepted) {
				return true
			}
		}
	}

	return false
}

func parseAcceptHeader(acceptHeader string) []string {
	var accepted []string
	for _, a := range strings.Split(acceptHeader, ",") {
		mediaType, _, err := mime.ParseMediaType(a)
		if err == nil {
			accepted = append(accepted, mediaType)
		}
	}
	return accepted
}

func containsWildcard(mediaTypes ...string) bool {
	for _, mediaType := range mediaTypes {
		if mediaType == allMediaTypes {
			return true
		}
	}
	return false
}
