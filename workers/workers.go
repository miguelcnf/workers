package workers

import (
	"log"
	"time"
)

type Pool interface {
	Run(func(...interface{}) (interface{}, error), []interface{}) chan TaskResult
}

type task struct {
	function func(...interface{}) (interface{}, error)
	args     []interface{}
	result   chan TaskResult
}

type TaskResult struct {
	Result interface{}
	Error  error
}

type Workers struct {
	taskCh chan task
}

func NewWorkerPool(maxConcurrency int) *Workers {
	taskCh := make(chan task, maxConcurrency)

	for i := 0; i < maxConcurrency; i++ {
		go func() {
			for {
				t := <-taskCh

				taskResult := TaskResult{}
				res, err := t.function(t.args...)
				if err != nil {
					taskResult.Error = err
				}

				taskResult.Result = res

				t.result <- taskResult
			}
		}()
	}

	return &Workers{
		taskCh: taskCh,
	}
}

func (w *Workers) Run(function func(args ...interface{}) (result interface{}, err error), args []interface{}) chan TaskResult {
	resultCh := make(chan TaskResult)

	t := task{
		function: function,
		args:     args,
		result:   resultCh,
	}

	for {
		select {
		case w.taskCh <- t:
			return resultCh
		default:
			log.Print("info: blocked waiting for available goroutines")

			time.Sleep(10 * time.Millisecond)
		}
	}
}
