package tracers

import (
	"go.opentelemetry.io/otel/trace"
)

// Span defines the span of a trace.
type Span struct {
	Name string
	span trace.Span
}

func newSpan(name string, span trace.Span) Span {
	newSpan := Span{
		Name: name,
		span: span,
	}

	return newSpan
}

// End ends the given span.
func (s *Span) End() {
	s.span.End()
}
