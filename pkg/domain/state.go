package domain

type State int

func (s State) Is(v State) bool { return s == v }

const (
	StateDraft State = iota
	StateCreated
	StateDeleted
)
