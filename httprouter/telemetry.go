package httprouter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NotifyErr notifies a tracer of an error that occurred while processing a request.
func NotifyErr(r *http.Request, err error, statusCode int) {
	tracer := trace.SpanFromContext(r.Context())

	if tracer != nil {
		tracer.AddEvent("error", trace.WithAttributes(
			attribute.String("uri", r.RequestURI),
			attribute.Int("statusCode", statusCode),
		))

		for k, v := range getKeyValueParams(r) {
			tracer.AddEvent("param", trace.WithAttributes(
				attribute.String(k, v),
			))
		}

		tracer.RecordError(err)
	}
}

func getKeyValueParams(r *http.Request) map[string]string {
	urlParams := chi.RouteContext(r.Context()).URLParams

	params := make(map[string]string, len(urlParams.Keys))

	for index, key := range urlParams.Keys {
		value := urlParams.Values[index]
		params[key] = value
	}

	return params
}
