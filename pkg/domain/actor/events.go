package actor

const CreatedEventType = "actor.created"

type CreatedEventData struct {
	Name string
	X, Y int
}

const MovedEventType = "actor.moved"

type MovedEventData struct {
	X, Y int
}
