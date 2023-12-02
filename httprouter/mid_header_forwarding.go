package httprouter

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// HeaderForwarder decorates a request context with the value of certain headers
// in order to allow transport.HTTPRequester to use those headers in outgoing requests.
func HeaderForwarder(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			_, span := tracer.Start(ctx, "HeaderForwarder")
			defer span.End()

			propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			span.AddEvent("HeaderForwarder processing")

			r2 := r.WithContext(ctx)
			next.ServeHTTP(w, r2)
		})
	}
}
