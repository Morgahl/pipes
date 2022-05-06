package pipes

func FanOut[T any](count, size int, in <-chan T) []ChanPull[T] {
	outs := make([]ChanPull[T], count)
	fan := make([]chan<- T, count)
	for i := range outs {
		ch := make(chan T, size)
		outs[i] = ch
		fan[i] = ch
	}

	go fanOutWorker(fan, in)

	return outs
}

func fanOutWorker[T any](fan []chan<- T, in <-chan T) {
	defer func() {
		for _, out := range fan {
			close(out)
		}
	}()

	for t := range in {
		for _, out := range fan {
			out <- t
		}
	}
}
