package fruits

import (
	"context"
)

// MonitorCounter request counter.
type MonitorCounter interface {
	CountRequest()
	CountSuccess()
	CountError()
}

type FruitMiddleware struct {
	next    *Service
	counter MonitorCounter
}

func NewFruitMiddleware(service *Service, counter MonitorCounter) *FruitMiddleware {
	fruitMiddleware := FruitMiddleware{
		next:    service,
		counter: counter,
	}

	return &fruitMiddleware
}

// GetFruitWithID get the fruit with the given id.
func (w *FruitMiddleware) GetFruitWithID(ctx context.Context, fruitID int64) (*Fruit, error) {
	w.counter.CountRequest()

	fruit, err := w.next.GetFruitWithID(ctx, fruitID)
	if err != nil {
		w.counter.CountError()

		return fruit, err
	}

	w.counter.CountSuccess()

	return fruit, err
}

// Create creates a fruit.
func (w *FruitMiddleware) Create(ctx context.Context, newfruit NewFruit) (int64, error) {
	w.counter.CountRequest()

	fruitID, err := w.next.Create(ctx, newfruit)
	if err != nil {
		w.counter.CountError()

		return fruitID, err
	}

	w.counter.CountSuccess()

	return fruitID, err
}

// SearchFruits search fruits who match the given filters.
func (w *FruitMiddleware) SearchFruits(ctx context.Context, givenFilter SearchFruitFilter) (*SearchFruitsResult, error) {
	w.counter.CountRequest()

	result, err := w.next.SearchFruits(ctx, givenFilter)
	if err != nil {
		w.counter.CountError()

		return result, err
	}

	w.counter.CountSuccess()

	return result, err
}

// DatasetStatus check the status of the fruit dataset.
func (w *FruitMiddleware) DatasetStatus(ctx context.Context) DatasetStatus {
	w.counter.CountRequest()
	status := w.next.DatasetStatus(ctx)
	w.counter.CountSuccess()

	return status
}
