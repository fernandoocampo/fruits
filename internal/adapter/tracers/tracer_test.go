package tracers_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/tracers"
	"github.com/stretchr/testify/assert"
)

func TestTraceSpan(t *testing.T) {
	// Given
	expectedSpanOneName := "doSomethingOne"
	expectedSpanTwoName := "doSomethingTwo"
	traceName := "fruits"
	aPackage := packageMock{
		tracer: tracers.New(traceName),
	}
	// When
	aPackage.doSomethingOne(context.TODO())
	// Then
	assert.Equal(t, expectedSpanOneName, aPackage.spanOne.Name)
	assert.Equal(t, expectedSpanTwoName, aPackage.spanTwo.Name)
}

type packageMock struct {
	tracer  *tracers.Tracer
	spanOne *tracers.Span
	spanTwo *tracers.Span
}

func (p *packageMock) doSomethingOne(ctx context.Context) {
	newCtx, spanOne := p.tracer.Start(ctx, "doSomethingOne")
	p.doSomethingTwo(newCtx)
	spanOne.End()
	p.spanOne = &spanOne
}

func (p *packageMock) doSomethingTwo(ctx context.Context) {
	_, spanTwo := p.tracer.Start(ctx, "doSomethingTwo")
	spanTwo.End()
	p.spanTwo = &spanTwo
}
