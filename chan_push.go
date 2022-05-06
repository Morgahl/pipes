package pipes

type ChanPush[T any] chan<- T

// Close closes the channel. Any attempts to push to a closed channel will panic. Closing an already
// closed channel will return immeadiately.
func (c ChanPush[T]) Close() {
	close(c)
}

// TryPush is a non-blocking operation that attempts to push a T onto the channel. This returns true
// if the T was successfully pushed, false if the channel was blocked. It is exceedingly unlikely
// that you will ever successfully push onto an unbuffered channel.
func (c ChanPush[T]) TryPush(t T) (ok bool) {
	select {
	case c <- t:
		return true
	default:
		return false
	}
}

// Push is a blocking operation that pushes a T onto the channel. This blocks while the channel is
// full. This will panic if the channel is closed.
func (c ChanPush[T]) Push(t T) {
	c <- t
}
