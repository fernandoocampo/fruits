package tracers

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// Service contains data related to this service.
type Service struct {
	Name, Version string
}

// Tracer implements logic to do traceability.
type Tracer struct {
	name string
}

func New(name string) *Tracer {
	newTracer := Tracer{
		name: name,
	}

	return &newTracer
}

func newResource(service Service) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service.Name),
			semconv.ServiceVersionKey.String(service.Version),
			attribute.String("environment", "demo"),
		),
	)

	return r
}

// Start tracer starts a new span and add values to the context.
func (t *Tracer) Start(ctx context.Context, spanName string) (context.Context, Span) {
	newCtx, span := otel.Tracer(t.name).Start(ctx, spanName)

	newSpan := newSpan(spanName, span)

	return newCtx, newSpan
}
