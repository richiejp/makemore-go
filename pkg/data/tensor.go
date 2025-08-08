package data

import (
	"math"

	"github.com/gomlx/gomlx/backends"
	_ "github.com/gomlx/gomlx/backends/simplego"
	. "github.com/gomlx/gomlx/graph"
	"github.com/gomlx/gomlx/types/shapes"
	"github.com/gomlx/gomlx/types/tensors"
	// "github.com/gomlx/gopjrt/dtypes"
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
		// cs = ConvertDType(cs, dtypes.F16)
		
		sums := ReduceAndKeep(cs, ReduceSum, 1)
		probs := Div(cs, sums)

		return CumSum(probs, 1)
	})

	return e.Call(cst)[0]
} 

func SampleTensor(b backends.Backend, ps *tensors.Tensor, seed int64) *tensors.Tensor {
	e := NewExec(b, func(ps *Node) *Node {
		
		rngState := Const(ps.Graph(), RngStateFromSeed(seed))
		rngState, rns := RandomUniform(rngState, shapes.Make(ps.DType(), 1))

		rns = BroadcastToShape(rns, ps.Shape())
		bools := Sign(Sub(ps, rns))

		return ArgMax(bools, 1)
	})

	return e.Call(ps)[0]
}

func MyRandomNormal(rngState *Node, shape shapes.Shape) (newRngState, values *Node) {
	g := rngState.Graph()
	var u1, u2 *Node
	newRngState, u1 = RandomUniform(rngState, shape)
	u1 = Max(u1, Const(g, shape.DType.SmallestNonZeroValueForDType()))

	// newRngState = Const(g, RngStateFromSeed(42))
	newRngState, u2 = RandomUniform(newRngState, shape)

	values = Mul(
		Sqrt(MulScalar(Log(u1), -2)),
		Cos(MulScalar(u2, 2*math.Pi)))
	return
}

