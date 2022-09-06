package main

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Println(job)
				return
			}
			results <- job.execute(ctx)

		case <-ctx.Done():
			fmt.Printf("Cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool) Results(ctx context.Context) {
	go func() {
		for {
			select {
			case res, ok := <-wp.results:
				if !ok {
					return
				}
				fmt.Println(res)

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}

	close(wp.jobs)
}
