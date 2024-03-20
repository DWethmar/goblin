package domain

import "testing"

func TestState_Is(t *testing.T) {
	type args struct {
		v State
	}
	tests := []struct {
		name string
		s    State
		args args
		want bool
	}{
		{"test1", StateDraft, args{StateDraft}, true},
		{"test2", StateDraft, args{StateCreated}, false},
		{"test3", StateDraft, args{StateDeleted}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Is(tt.args.v); got != tt.want {
				t.Errorf("State.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
