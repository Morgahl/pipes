package pipes

import "sync"

// FanIn is a non-blocking operation that creates len(ins) goroutines and forwards each T read onto
// the returned push only channel of specified size. Each goroutine will exit after it's assigned
// input channel is closed and emptied. The last goroutine will close the returned channel to signal
// completion of work.
func FanIn[T any](size int, ins ...<-chan T) ChanPull[T] {
	out := make(chan T, size)
	if len(ins) < 1 {
		close(out)
		return out
	}

	go fanInCoordinator(ins, out)

	return out
}

func fanInCoordinator[T any](ins []<-chan T, out chan<- T) {
	defer close(out)

	wg := &sync.WaitGroup{}
	wg.Add(len(ins))
	for _, in := range ins[1:] {
		go fanInWorker(wg, in, out)
	}

	// demote to a worker to guarantee there is always one worker running and launch one less
	// goroutine
	fanInWorker(wg, ins[0], out)

	wg.Wait()
}

func fanInWorker[T any](wg *sync.WaitGroup, in <-chan T, out chan<- T) {
	defer wg.Done()

	for t := range in {
		out <- t
	}
}
