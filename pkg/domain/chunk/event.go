package chunk

const (
	CreatedEventType = "chunk.created"
)

type CreatedEventData struct {
	X, Y int
}
