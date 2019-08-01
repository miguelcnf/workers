package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"workers/workers"
)

const (
	maxConcurrency = 100
)

func main() {
	workerPool := workers.NewWorkerPool(maxConcurrency)

	for {
		x := rand.Intn(10-1) + 1
		y := rand.Intn(10-1) + 1

		var input []interface{}
		input = append(input, x)
		input = append(input, y)

		result := workerPool.Run(func(args ...interface{}) (result interface{}, err error) {

			// Emulate random slow tasks
			sleep := time.Duration(rand.Intn(500-10)+10) * time.Millisecond
			time.Sleep(sleep)

			i := 0
			for _, v := range args {
				val, ok := v.(int)
				if !ok {
					return nil, errors.New("invalid argument")
				}
				i += val
			}

			return i, nil
		}, input)

		go func() {
			res := <-result
			if res.Error != nil {
				log.Printf("error: %v\n", res.Error)
				os.Exit(1)
			}

			if sum, ok := res.Result.(int); ok {
				log.Printf("result: x(%v) + y(%v) = %v\n", x, y, sum)
			} else {
				log.Printf("error: unexpected result type")
			}
		}()
	}
}
