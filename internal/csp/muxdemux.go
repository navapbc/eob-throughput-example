package csp

// func Mux[T any](in chan []T, out []chan T) {
// 	next := 0
// 	for {
// 		p := <-in
// 		out[next] <- p
// 		next += 1
// 		if next == len(out) {
// 			next = 0
// 		}
// 	}
// }

// Technically, we don't need demux, because workers
// can write to a shared channel. This is an identity
// process, therefore, and could be omitted.
func Demux[T any](in chan T, out chan T) {
	for {
		v := <-in
		out <- v
	}
}
