package workers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	function = func(args ...interface{}) (result interface{}, err error) {
		i := 0
		for _, v := range args {
			val, ok := v.(int)
			if !ok {
				return nil, errors.New("invalid argument")
			}
			i += val
		}

		return i, nil
	}
)

func TestNewWorkerPool(t *testing.T) {
	maxConcurrency := 2
	workerPool := NewWorkerPool(maxConcurrency)

	assert.NotNil(t, workerPool)
}

func TestWorkers_Run(t *testing.T) {
	maxConcurrency := 2
	workerPool := NewWorkerPool(maxConcurrency)

	t.Run("should execute function and return result", func(t *testing.T) {
		var input []interface{}
		input = append(input, 1)
		input = append(input, 1)

		result := workerPool.Run(function, input)

		res := <-result
		require.NotNil(t, res)
		assert.NoError(t, res.Error)
		assert.Equal(t, 2, res.Result)
	})

	t.Run("should execute function and return error", func(t *testing.T) {
		var input []interface{}
		input = append(input, 1)
		input = append(input, "invalid")

		result := workerPool.Run(function, input)

		res := <-result
		require.NotNil(t, res)
		assert.Error(t, res.Error)
		assert.Equal(t, nil, res.Result)
	})
}
