package conc

import (
	"testing"
)

func TestGroupTasksByID(t *testing.T) {
	t.Run("should group tasks by id", func(t *testing.T) {
		in := make(chan Task[int], 2)

		out, err := GroupTasks(in, 2)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}

		in <- Task[int]{Group: 0, Value: 1}
		in <- Task[int]{Group: 1, Value: 2}

		close(in)

		got1 := <-out[0]
		want1 := 1
		if got1 != want1 {
			t.Errorf("got %v, want %v", got1, want1)
		}

		got2 := <-out[1]
		want2 := 2
		if got2 != want2 {
			t.Errorf("got %v, want %v", got2, want2)
		}
	})

	t.Run("should return an error when n is less than or equal to 0", func(t *testing.T) {
		in := make(chan Task[int], 2)
		_, err := GroupTasks(in, 0)
		if err == nil {
			t.Errorf("got nil, want an error")
		}
	})
}
