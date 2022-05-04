package pipes

import "sync"

func FanIn[T any](size int, ins ...<-chan T) ChanPull[T] {
	if len(ins) < 1 {
		return nil
	}
	out := make(chan T, size)
	go fanInCoordinator(ins, out)
	return out
}

func FanInExisting[T any](out chan<- T, ins ...<-chan T) {
	if len(ins) < 1 {
		return
	}
	go fanInCoordinator(ins, out)
	return
}

func fanInCoordinator[T any](ins []<-chan T, out chan<- T) {
	defer close(out)
	wg := &sync.WaitGroup{}
	wg.Add(len(ins))
	for _, in := range ins[1:] {
		go fanInWorker(wg, in, out)
	}
	// demote to a worker to guarantee there is always one worker running and launch one less goroutine
	fanInWorker(wg, ins[0], out)
	wg.Wait()
}

func fanInWorker[T any](wg *sync.WaitGroup, in <-chan T, out chan<- T) {
	defer wg.Done()
	for t := range in {
		out <- t
	}
}
