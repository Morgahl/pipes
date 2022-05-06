package pipes

func Router[T any, N comparable](size int, matches []N, compare func(T) N, in <-chan T) ([]ChanPull[T], ChanPull[T]) {
	orElse := make(chan T, size)
	outs := make([]ChanPull[T], len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i, match := range matches {
		out := make(chan T, size)
		outs[i] = out
		routes[match] = out
	}

	go routerWorker(compare, in, routes, orElse)

	return outs, orElse
}

func routerWorker[T any, N comparable](compare func(T) N, in <-chan T, routes map[N]chan<- T, orElse chan<- T) {
	defer func() {
		for _, out := range routes {
			close(out)
		}

		close(orElse)
	}()

	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
			continue
		}

		orElse <- t
	}
}

func RouterWithSink[T any, N comparable](size int, matches []N, compare func(T) N, sink func(T), in <-chan T) []ChanPull[T] {
	outs := make([]ChanPull[T], len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i, match := range matches {
		out := make(chan T, size)
		outs[i] = out
		routes[match] = out
	}

	go routerWithSinkWorker(compare, in, routes, sink)

	return outs
}

func routerWithSinkWorker[T any, N comparable](compare func(T) N, in <-chan T, routes map[N]chan<- T, sink func(T)) {
	defer func() {
		for _, out := range routes {
			close(out)
		}
	}()

	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
			continue
		}

		sink(t)
	}
}

func RoundRobin[T any](size, count int, in <-chan T) []ChanPull[T] {
	if count < 1 {
		return nil
	}

	return Distribute(size, count, roundRobinChooser[T](count), in)
}

func roundRobinChooser[T any](count int) func(T) int {
	lastIdx := 0
	return func(t T) int {
		if lastIdx >= count {
			lastIdx = 0
		}

		idx := lastIdx
		lastIdx++

		return idx
	}
}

func Distribute[T any](size, count int, choose func(T) int, in <-chan T) []ChanPull[T] {
	if count < 1 {
		return nil
	}

	outs := make([]ChanPull[T], count)
	pushes := make([]chan<- T, count)
	for i := 0; i < count; i++ {
		ch := make(chan T, size)
		outs[i] = ch
		pushes[i] = ch
	}

	go distrbuteWorker(choose, in, pushes)

	return outs
}

func distrbuteWorker[T any](choose func(T) int, in <-chan T, outs []chan<- T) {
	defer func() {
		for _, out := range outs {
			close(out)
		}
	}()

	for t := range in {
		outs[choose(t)] <- t
	}
}
