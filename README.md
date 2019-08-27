# Normal Distribution Z-Score Table - Go

Generates a normal distribution z-score table in Go. Provides functions for quick z-score and percentage lookup.

## Installing

To start using ztable, install Go and run `go get`:

```sh
$ go get -u github.com/gregscott94/z-table-golang
```

## To Use

```go
// Creating a new z-score table
zTable := NewZTable(nil)

// To find the percentage of a given z-score
percentage := zTable.FindPercentage(1.09)
// percentage = 0.8621434279679557

// To find the closest z-score given a percentage
zScore, err := zTable.FindZScore(0.04363)
// zScore = -1.71
```
