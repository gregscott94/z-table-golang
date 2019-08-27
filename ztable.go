package main

import (
	"log"
	"math"
)

type ZTable struct {
	zScoreMap   map[float64]int
	percentages []float64
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
		zScoreMap:   make(map[float64]int),
		percentages: []float64{},
	}
	zScore := float64(-4)
	leafNodes := []*Node{}
	for zScore <= 4 {
		percentage := (1 / math.Sqrt(2*math.Pi)) * math.Exp((-1*zScore*zScore)/2)
		index := len(zTable.percentages)
		if index%bucketSize == 0 {
			leafNodes = append(leafNodes, &Node{
				value: percentage,
				index: index,
			})
		}
		zTable.zScoreMap[zScore] = index
		zTable.percentages = append(zTable.percentages, percentage)
		zScore = zScore + 0.01
	}
	return &zTable
}

func main() {
	log.Println("Init commit")
}
