package httprouter

import (
	"net/http"

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

		for k, v := range Params(r) {
			tracer.AddEvent("param", trace.WithAttributes(
				attribute.String(k, v),
			))
		}

		tracer.RecordError(err)
	}
}
