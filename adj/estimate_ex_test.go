package adj_test

import (
	"fmt"

	"github.com/soniakeys/graph2"
	"github.com/soniakeys/graph2/adj"
)

// Define a type to use as Node.Data.
type est struct{}

// Implement graph2.Estimator.
func (est) Estimate(graph2.EstimateNode) float64 {
	return 4
}

func ExampleNode_Estimate() {
	// Use the type that implents graph2.Estimator as the Data field.
	// Use adj.Node as a graph2.EstimateNode.
	var n graph2.EstimateNode = &adj.Node{Data: est{}}
	// Node.Estimate will call Data.Estimate.
	fmt.Println(n.Estimate(nil))
	// Output:
	// 4
}
