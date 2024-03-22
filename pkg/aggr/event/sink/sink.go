package sink

import (
	"context"
	"hash/fnv"
	"sync"

	"github.com/dwethmar/goblin/pkg/aggr"
)

type WorkerSink struct {
	workerChannels []chan *aggr.Event
	numWorkers     int
	handlers       []aggr.EventHandler
	wg             sync.WaitGroup
	errCh          chan error
}

func NewWorkerSink(handlers []aggr.EventHandler) *WorkerSink {
	numWorkers := len(handlers)

	ws := &WorkerSink{
		numWorkers:     numWorkers,
		workerChannels: make([]chan *aggr.Event, numWorkers),
		handlers:       handlers,
		errCh:          make(chan error),
	}

	for i := 0; i < numWorkers; i++ {
		ws.workerChannels[i] = make(chan *aggr.Event, 100) // Adjust buffer size as needed
		ws.wg.Add(1)
		go ws.worker(i)
	}

	return ws
}

func (ws *WorkerSink) worker(workerID int) {
	defer ws.wg.Done()
	for event := range ws.workerChannels[workerID] {
		if err := ws.handlers[workerID].HandleEvent(context.Background(), event); err != nil {
			ws.errCh <- err
		}
	}
}

func (ws *WorkerSink) ProcessEvent(ctx context.Context, event *aggr.Event) {
	index := ws.getWorkerIndex(event.AggregateID)
	ws.workerChannels[index] <- event
}

func (ws *WorkerSink) getWorkerIndex(aggregateID string) int {
	h := fnv.New32a()
	h.Write([]byte(aggregateID))
	return int(h.Sum32()) % ws.numWorkers
}

func (ws *WorkerSink) Errors() <-chan error {
	return ws.errCh
}

func (ws *WorkerSink) Close() {
	for _, channel := range ws.workerChannels {
		close(channel)
	}
	ws.wg.Wait()
}
