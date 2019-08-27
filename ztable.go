package ztable

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"strconv"

	"gonum.org/v1/gonum/integrate/quad"
)

type ZTable struct {
	zScoreMap map[string]int
	leafNodes []*LeafNode
	rootNode  *Node
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

type LeafNode struct {
	zScore     string
	percentage float64
}

func (zt *ZTable) FindPercentage(zScore float64) float64 {
	key := fmt.Sprintf(`%.2f`, zScore)
	if index, ok := zt.zScoreMap[key]; ok {
		return zt.leafNodes[index].percentage
	}
	return 0
}

func (zt *ZTable) FindZScore(percentage float64) (float64, error) {
	currNode := zt.rootNode
	startingIndex := 0
	if currNode != nil {
		for {
			if percentage > currNode.value && currNode.rightNode != nil {
				currNode = currNode.rightNode
			} else if percentage <= currNode.value && currNode.leftNode != nil {
				currNode = currNode.leftNode
			} else {
				startingIndex = currNode.index
				break
			}
		}
		if startingIndex == 0 && percentage < zt.leafNodes[0].percentage {
			return strconv.ParseFloat(zt.leafNodes[0].zScore, 64)
		}
		for startingIndex < len(zt.leafNodes) {
			currLeaf := zt.leafNodes[startingIndex]
			if percentage == currLeaf.percentage || startingIndex+1 >= len(zt.leafNodes) {
				return strconv.ParseFloat(currLeaf.zScore, 64)
			} else if nextLeaf := zt.leafNodes[startingIndex+1]; percentage < nextLeaf.percentage {
				if percentage-currLeaf.percentage <= nextLeaf.percentage-percentage {
					return strconv.ParseFloat(currLeaf.zScore, 64)
				}
				return strconv.ParseFloat(nextLeaf.zScore, 64)
			}
			startingIndex++
		}
	}
	return 0, errors.New("Unable to find ZScore given percentage")
}

func NewZTable(options *ZTableOptions) *ZTable {
	bucketSize := 30
	if options != nil && options.BucketSize != 0 {
		bucketSize = options.BucketSize
	}
	zTable := ZTable{
		zScoreMap: make(map[string]int),
		leafNodes: []*LeafNode{},
	}

	zScore := float64(-4)
	for zScore <= 4 {
		concurrent := runtime.GOMAXPROCS(0)
		percentage := quad.Fixed(normalProbabilityDensity, math.Inf(-1), zScore, 1000, nil, concurrent)
		index := len(zTable.leafNodes)
		if index%bucketSize == 0 {

		}
		zScoreString := fmt.Sprintf(`%.2f`, zScore)
		zTable.zScoreMap[zScoreString] = index
		zTable.leafNodes = append(zTable.leafNodes, &LeafNode{zScore: zScoreString, percentage: percentage})
		zScore = zScore + 0.01
	}

	initLayer := []*Node{}
	i := 0
	for i < len(zTable.leafNodes) {
		if i+bucketSize < len(zTable.leafNodes) {
			initLayer = append(initLayer, &Node{
				value: zTable.leafNodes[i+bucketSize-1].percentage,
				index: i,
			})
		}
		i = i + bucketSize
	}

	rootNodeSlice := buildTree(initLayer)
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
