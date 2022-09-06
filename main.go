package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {

	wp := New(runtime.NumCPU())
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*200)
	defer cancel()
	xJob := TaskGenerator(ctx)
	go wp.GenerateFrom(xJob)
	wp.Run(ctx)
	fmt.Println(wp.Results())

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
		}

		xJob = append(xJob, job)

	}

	return xJob

}
