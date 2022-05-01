package async

import (
	"sync"

	"github.com/curlymon/pipes/fn"
)

// really the big thing that isn't obvious here is you lose any ordering going through
// damn near everything in the package as planned lol.
func Map[T any, N any](count, size int, mp fn.Map[T, N], in <-chan T) <-chan N {
	out := make(chan N, size)
	go mapCoordinator(count, mp, in, out)
	return out
}

func mapCoordinator[T any, N any](count int, mp fn.Map[T, N], in <-chan T, out chan<- N) {
	if count < 1 {
		count = 1
	}
	defer close(out)
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for ; count > 1; count-- {
		go mapWorker(wg, mp, in, out)
	}
	// demote to a worker to guarantee there is always one worker running and launch one less goroutine
	mapWorker(wg, mp, in, out)
	wg.Wait()
}

func mapWorker[T any, N any](wg *sync.WaitGroup, mp fn.Map[T, N], in <-chan T, out chan<- N) {
	defer wg.Done()
	for t := range in {
		out <- mp(t)
	}
}

func MapWithError[T any, N any](count, size int, mp fn.MapWithError[T, N], in <-chan T) (<-chan N, <-chan error) {
	out, err := make(chan N, size), make(chan error, size)
	go mapWithErrorCoordinator(count, mp, in, out, err)
	return out, err
}

func mapWithErrorCoordinator[T any, N any](count int, mp fn.MapWithError[T, N], in <-chan T, out chan<- N, err chan<- error) {
	if count < 1 {
		count = 1
	}
	defer func() { close(out); close(err) }()
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for ; count > 1; count-- {
		go mapWithErrorWorker(wg, mp, in, out, err)
	}
	// demote to a worker to guarantee there is always one worker running and launch one less goroutine
	mapWithErrorWorker(wg, mp, in, out, err)
	wg.Wait()
}

func mapWithErrorWorker[T any, N any](wg *sync.WaitGroup, mp fn.MapWithError[T, N], in <-chan T, out chan<- N, err chan<- error) {
	defer wg.Done()
	for t := range in {
		if n, er := mp(t); er != nil {
			err <- er
		} else {
			out <- n
		}
	}
}

func MapWithErrorSink[T any, N any](count, size int, mp fn.MapWithError[T, N], sink fn.Sink[error], in <-chan T) <-chan N {
	out := make(chan N, size)
	go mapWithErrorSinkCoordinator(count, mp, sink, in, out)
	return out
}

func mapWithErrorSinkCoordinator[T any, N any](count int, mp fn.MapWithError[T, N], sink fn.Sink[error], in <-chan T, out chan<- N) {
	if count < 1 {
		count = 1
	}
	defer close(out)
	wg := &sync.WaitGroup{}
	wg.Add(count)
	for ; count > 1; count-- {
		go mapWithErrorSinkWorker(wg, mp, sink, in, out)
	}
	// demote to a worker to guarantee there is always one worker running and launch only `count` goroutines
	mapWithErrorSinkWorker(wg, mp, sink, in, out)
	wg.Wait()
}

func mapWithErrorSinkWorker[T any, N any](wg *sync.WaitGroup, mp fn.MapWithError[T, N], sink fn.Sink[error], in <-chan T, out chan<- N) {
	defer wg.Done()
	for t := range in {
		if n, er := mp(t); er != nil {
			sink(er)
		} else {
			out <- n
		}
	}
}
