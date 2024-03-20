package actor

const (
	CreatedEventType   = "actor.created"
	DestroyedEventType = "actor.destroyed"
	MovedEventType     = "actor.moved"
)

type CreatedEventData struct {
	Name string
	X, Y int
}

type MovedEventData struct {
	X, Y int
}
