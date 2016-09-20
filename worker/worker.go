package worker

// Job a job to be added to a worker needs implements this interface
// type Job func(i interface{}) error
type Job interface {
	Exec() error
}

// NewWorker creates a worker
func NewWorker(id int, workerQueue chan chan Job) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan Job),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Worker TODO
type Worker struct {
	ID          int
	Work        chan Job
	WorkerQueue chan chan Job
	QuitChan    chan bool
}

// Start starts the worker
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				work.Exec()

			case <-w.QuitChan:
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
