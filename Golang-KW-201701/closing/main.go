package main

func main() {
	in := make(chan bool)

	// START1 OMIT
	for msg := range in {
		// ...
		_ = msg // OMIT
	}
	// STOP1 OMIT

	// START2 OMIT
	select {
	case msg := <-in:
		// If msg is zero, return
		_ = msg // OMIT
	}
	// STOP2 OMIT

	// START3 OMIT
	select {
	case msg, running := <-in:
		if !running {
			return
		}
		_ = msg // OMIT
	}
	// STOP3 OMIT
}
