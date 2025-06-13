package processing

import (
	"log/slog"
	"time"

	"insider-egemen-avci/backend-path-p1/internal/models"
)

type Pool struct {
	numberOfWorkers int
	jobs            chan models.Transaction
	results         chan models.Transaction
}

func NewPool(numberOfWorkers, queueSize int) *Pool {
	return &Pool{
		numberOfWorkers: numberOfWorkers,
		jobs:            make(chan models.Transaction, queueSize),
		results:         make(chan models.Transaction, queueSize),
	}
}

func (pool *Pool) worker(id int) {
	for job := range pool.jobs {
		slog.Info("Worker", "id", id, "processing transaction", "transaction_id", job.ID)

		time.Sleep(time.Second * 2)
		job.Complete()

		slog.Info("Worker", "id", id, "completed transaction", "transaction_id", job.ID)
		pool.results <- job
	}
}

func (pool *Pool) Start() {
	slog.Info("Starting pool with", "numberOfWorkers", pool.numberOfWorkers)

	for i := 0; i < pool.numberOfWorkers; i++ {
		go pool.worker(i)
	}
}

func (pool *Pool) AddJob(job models.Transaction) {
	pool.jobs <- job
}

func (pool *Pool) CloseJobs() {
	close(pool.jobs)
}

func (pool *Pool) Results() <-chan models.Transaction {
	return pool.results
}
