package httprouter

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// HeaderForwarder decorates a request context with the value of certain headers
// in order to allow transport.HTTPRequester to use those headers in outgoing requests.
func HeaderForwarder(tracer trace.Tracer) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			_, span := tracer.Start(ctx, "HeaderForwarder")
			defer span.End()

			propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			span.AddEvent("HeaderForwarder processing")

			r2 := r.WithContext(ctx)
			handler(w, r2)
		}
	}
}
