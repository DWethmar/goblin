package conc

import "fmt"

type Task[T any] struct {
	Group int // Group is used to determine which channel to send the task to.
	Value T
}

// GroupTasks makes sure that tasks with the same group are sent to the same channel.
func GroupTasks[T any](in <-chan Task[T], n int) ([]chan T, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

	out := make([]chan T, n)
	for i := 0; i < n; i++ {
		out[i] = make(chan T)
	}

	go func() {
		for task := range in {
			i := task.Group % len(out)
			out[i] <- task.Value
		}

		// Close all output channels
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	return out, nil
}
