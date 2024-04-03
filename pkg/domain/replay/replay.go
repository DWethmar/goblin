// Package replay provides a way to replay events from the event store to the event bus.
package replay

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
	"github.com/dwethmar/goblin/pkg/aggr/sink"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	actormemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/domain/chunk"
	chunkmemory "github.com/dwethmar/goblin/pkg/domain/chunk/memory"
)

type Replayer struct {
	logger     *slog.Logger
	eventStore event.Repository
	eventBus   *aggr.EventBus
}

func New(
	logger *slog.Logger,
	eventStore event.Repository,
	eventBus *aggr.EventBus,
) *Replayer {
	return &Replayer{
		logger:     logger,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

// Replay all events from the event store and rebuild the whole state.
func (r *Replayer) Replay(
	ctx context.Context,
	actorRepo actor.Repository,
	chunkRepo chunk.Repository,
) error {
	actorRepos := []actor.Repository{}
	chunkRepos := []chunk.Repository{}

	handlers := []aggr.EventHandler{}
	for range 1 {
		ar := actormemory.NewRepository()
		actorRepos = append(actorRepos, ar)

		cr := chunkmemory.NewRepository()
		chunkRepos = append(chunkRepos, cr)

		eventBus := aggr.NewEventBus()
		eventBus.Subscribe(actor.MatchAllEvents, actor.ActorSinkHandler(ar))
		eventBus.Subscribe(chunk.MatchAllEvents, chunk.ChunkSinkHandler(cr))
		handlers = append(handlers, eventBus)
	}

	errCn := make(chan error)
	eventSink := sink.New(r.eventStore.All(errCn), handlers)

	done := make(chan struct{})
	go func() {
		eventSink.Sink(ctx)
		done <- struct{}{}
	}()

	// wait for the event sink to be ready
Loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			break Loop
		case err := <-errCn:
			return err
		case err := <-eventSink.Errors():
			return err
		}
	}

	wg := sync.WaitGroup{}
	var errr error

	// write all actors to the main actor repo
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := actor.Merge(ctx, actorRepo, actorRepos...); err != nil {
			errr = fmt.Errorf("merge actors: %w", err)
			return
		}
		fmt.Printf("merged actors :D\n")
	}()

	// write all chunk repos to the main chunk repo
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := chunk.Merge(ctx, chunkRepo, chunkRepos...); err != nil {
			errr = fmt.Errorf("merge chunks: %w", err)
			return
		}
		fmt.Printf("merged chunks\n")
	}()

	wg.Wait()
	if errr != nil {
		return errr
	}

	fmt.Println("replay done")

	return nil
}
