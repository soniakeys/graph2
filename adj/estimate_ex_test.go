package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
)

// Define a type to use as Node.Data
type est struct{}

// The type must impliment graph.Estimator.
func (est) Estimate(graph.EstimateNode) float64 {
	return 4
}

func ExampleNode_Estimate() {
	// Use the type as the Data field.
	n := adj.Node{Data: est{}}
	// Node.Estimate will call Data.Estimate.
	fmt.Println(n.Estimate(nil))
	// Output:
	// 4
}
