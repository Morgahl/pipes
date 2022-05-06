package pipes

import "sync"

// FanIn is a non-blocking operation that creates len(ins) goroutines and forwards each T read onto
// the returned push only channel of specified size. Each goroutine will exit after it's assigned
// input channel is closed and emptied. The last goroutine will close the returned pull only channel
// to signal completion of processing.
func FanIn[T any](size int, ins ...<-chan T) ChanPull[T] {
	out := make(chan T, size)
	if len(ins) < 1 {
		// Let's never return a nil channel, a close empty channel has better behaviors
		close(out)
		return out
	}

	go fanInCoordinator(ins, out)

	return out
}

// fanInCoordinator will create len(ins)-1 fanInWorkers then demote itself to a fanInWorker, ensure
// that at least 1 push only channel is passed or the function will panic. This will close the
// passed push only channel after the last worker exits.
func fanInCoordinator[T any](ins []<-chan T, out chan<- T) {
	defer close(out)

	wg := &sync.WaitGroup{}
	wg.Add(len(ins))
	// skipping the first create a worker for each passed pull only channel
	for _, in := range ins[1:] {
		go fanInWorker(wg, in, out)
	}

	// demote to a worker to guarantee there is always one worker running and launch one less
	// goroutine
	fanInWorker(wg, ins[0], out)

	wg.Wait()
}

// fanInWorker iterates over the passed pull only channel forwarding values to the passed push only
// channel. Worker will exit after the pul only channel is closed and emptied.
func fanInWorker[T any](wg *sync.WaitGroup, in <-chan T, out chan<- T) {
	defer wg.Done()

	for t := range in {
		out <- t
	}
}
