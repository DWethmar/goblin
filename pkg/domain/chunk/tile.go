package chunk

// Tile is a single tile in a chunk
type Tile struct {
	X     int // X within the chunk
	Y     int // Y within the chunk
	Value int
}

// GlobalPosition returns the global position of the tile
func (t Tile) GlobalPosition(chunkX, chunkY, chunkWidth, chunkHeight int) (int, int) {
	return chunkX*chunkWidth + t.X, chunkY*chunkHeight + t.Y
}
