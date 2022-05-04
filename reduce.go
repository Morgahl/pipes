package pipes

import "time"

func Reduce[T any, Acc any](reduce func(T, Acc) Acc, acc Acc, in <-chan T) Acc {
	for t := range in {
		acc = reduce(t, acc)
	}
	return acc
}

func ReduceAndEmit[T any, Acc any](reduce func(T, Acc) Acc, acc Acc, in <-chan T) ChanPull[Acc] {
	// we only expect to emit a single value and then close the out chan immeadiately
	// after processing. This allows the goroutine to exit without forcing it to sync
	// with the recieving goroutine.
	out := make(chan Acc, 1)
	go reduceWorker(reduce, acc, in, out)
	return out
}

func reduceWorker[T any, Acc any](reduce func(T, Acc) Acc, acc Acc, in <-chan T, out chan<- Acc) {
	defer close(out)
	for t := range in {
		acc = reduce(t, acc)
	}
	out <- acc
}

// TODO: Decide if Window's reduce func should take a time.Time object as well and will be passed the "tick" from the
// ticker for use internally for structuring the Acc being emitted.
func Window[T any, Acc any](size int, window time.Duration, reduce func(T, Acc) Acc, acc func() Acc, in <-chan T) ChanPull[Acc] {
	out := make(chan Acc, 1)
	go windowWorker(window, reduce, acc, in, out)
	return out
}

func windowWorker[T any, Acc any](window time.Duration, reduce func(T, Acc) Acc, acc func() Acc, in <-chan T, out chan<- Acc) {
	defer close(out)
	ac := acc()
	ticker := time.NewTicker(window)
	defer ticker.Stop()
	for {
		select {
		case t, ok := <-in:
			if !ok {
				out <- ac
				return
			}
			ac = reduce(t, ac)
		case <-ticker.C:
			out <- ac
			ac = acc()
		}
	}
}
