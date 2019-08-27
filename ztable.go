package main

import (
	"fmt"
	"math"
	"runtime"

	"gonum.org/v1/gonum/integrate/quad"
)

type ZTable struct {
	zScoreMap   map[string]int
	percentages []float64
	rootNode    *Node
}

type ZTableOptions struct {
	BucketSize int
}

type Node struct {
	value     float64
	leftNode  *Node
	rightNode *Node
	index     int
}

func NewZTable(options *ZTableOptions) *ZTable {
	bucketSize := 10
	if options != nil && options.BucketSize != 0 {
		bucketSize = options.BucketSize
	}
	zTable := ZTable{
		zScoreMap:   make(map[string]int),
		percentages: []float64{},
	}
	zScore := float64(-4)
	leafNodes := []*Node{}
	for zScore <= 4 {
		concurrent := runtime.GOMAXPROCS(0)
		percentage := quad.Fixed(normalProbabilityDensity, math.Inf(-1), zScore, 1000, nil, concurrent)
		index := len(zTable.percentages)
		if index%bucketSize == 0 {
			leafNodes = append(leafNodes, &Node{
				value: percentage,
				index: index,
			})
		}
		zTable.zScoreMap[fmt.Sprintf(`%.2f`, zScore)] = index
		zTable.percentages = append(zTable.percentages, percentage)
		zScore = zScore + 0.01
	}
	rootNodeSlice := buildTree(leafNodes)
	if len(rootNodeSlice) > 0 {
		zTable.rootNode = rootNodeSlice[0]
	}
	return &zTable
}

func normalProbabilityDensity(x float64) float64 {
	return (1 / math.Sqrt(2*math.Pi)) * math.Exp((-1*x*x)/2)
}

func buildTree(layer []*Node) []*Node {
	currLayer := []*Node{}
	count := 0
	for count < len(layer)-1 {
		currLayer = append(currLayer, &Node{
			value:     (layer[count].value + layer[count+1].value) / 2,
			leftNode:  layer[count],
			rightNode: layer[count+1],
		})
		count = count + 1
	}
	if len(currLayer) <= 1 {
		return currLayer
	}
	return buildTree(currLayer)
}

func main() {
	zt := NewZTable(nil)
	fmt.Println(zt)
}
