package pipes

// TODO: implement FanIn handler.
// We do this as N-to-1 where there are N goroutines each listening to 1 of N channels
// and each forwarding messages into a single output channel.
