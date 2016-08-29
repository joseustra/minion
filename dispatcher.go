package minion

import "fmt"

// Dispatcher TODO
type Dispatcher struct {
	Size int
}

// NewDispatcher TODO
func NewDispatcher(size int) *Dispatcher {
	return &Dispatcher{Size: size}
}

// WorkerQueue TODO
var WorkerQueue chan chan Job

// WorkQueue A buffered channel that we can send work requests on.
var WorkQueue = make(chan Job, 100)

// StartDispatcher TODO
func (d *Dispatcher) StartDispatcher() {
	WorkerQueue = make(chan chan Job, d.Size)

	for i := 0; i < d.Size; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()
}

// Add adds a new job to the worker
func (d *Dispatcher) Add(job Job) {
	WorkQueue <- job
}
