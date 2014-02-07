package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/adj"
)

// Define a type to use as Node.Data.
type est struct{}

// Implement graph.Estimator.
func (est) Estimate(graph.EstimateNode) float64 {
	return 4
}

func ExampleNode_Estimate() {
	// Use the type that implents graph.Estimator as the Data field.
	// Use adj.Node as a graph.EstimateNode.
	var n graph.EstimateNode = &adj.Node{Data: est{}}
	// Node.Estimate will call Data.Estimate.
	fmt.Println(n.Estimate(nil))
	// Output:
	// 4
}
