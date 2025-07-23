package data

import (
	"github.com/gomlx/gomlx/backends"
	_ "github.com/gomlx/gomlx/backends/simplego"
	. "github.com/gomlx/gomlx/graph"
	"github.com/gomlx/gomlx/types/tensors"
)

func ToProbabilities(backend backends.Backend, counts [][]uint16) *tensors.Tensor {
	var cs [][]float32

	// GoMLX doesn't support creating tensors from uint16
	// Also drop the last row which is bigrams starting with <E>
	for _, row := range counts[:len(counts) - 1] {
		var r []float32

		for _, count := range row {
			r = append(r, float32(count))
		}
		cs = append(cs, r)
	}

	cst := tensors.FromValue(cs)

	e := NewExec(backend, func(cs *Node) *Node {
		// Not supported by simple Go backend
		//cs = ConvertDType(cs, dtypes.F16)
		
		sums := ReduceAndKeep(cs, ReduceSum, 1)

		return Div(cs, sums)
	})

	return e.Call(cst)[0]
} 

