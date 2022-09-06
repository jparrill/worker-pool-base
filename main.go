package main

import (
	"context"
	"errors"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	wp := New(runtime.NumCPU())
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*20)
	defer cancel()
	xJob := TaskGenerator(ctx)
	go wp.GenerateFrom(xJob)
	wg.Add(1)
	go wp.Results(ctx, &wg)
	wp.Run(ctx)
	wg.Wait()
}

func TaskGenerator(ctx context.Context) []Job {
	var jtype string
	errDefault := errors.New("wrong argument type")
	xJob := make([]Job, 0, 0)

	for i := 0; i <= 100; i++ {
		if i%2 == 0 {
			jtype = "even"
		} else {
			jtype = "odd"
		}

		id := strconv.Itoa(i)

		execFn := func(ctx context.Context, args interface{}) (interface{}, error) {
			argVal, ok := args.(int)
			if !ok {
				return nil, errDefault
			}

			return argVal * i, nil
		}

		jobD := JobDescriptor{
			ID:    JobID(id),
			JType: jobType(jtype),
			Metadata: jobMetadata{
				"ID": i,
			},
		}

		job := Job{
			Descriptor: jobD,
			ExecFn:     execFn,
			Args:       i,
		}

		xJob = append(xJob, job)

	}

	return xJob

}
