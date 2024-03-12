package game

import "fmt"

func (g *Game) CreateActor(id, name string) error {
	if err := g.actorService.CreateActor(id, name); err != nil {
		return fmt.Errorf("could not create actor: %w", err)
	}

	return nil
}
