package pipes

type ChanPush[T any] chan<- T

func (c ChanPush[T]) Close() {
	close(c)
}

func (c ChanPush[T]) TryPush(t T) (ok bool) {
	select {
	case c <- t:
		return true
	default:
		return false
	}
}

func (c ChanPush[T]) Push(t T) {
	c <- t
}

func (c ChanPush[T]) FanIn(size int, ins ...<-chan T) {
	FanInExisting(c, ins...)
}
