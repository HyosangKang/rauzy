package rauzy

import (
	"testing"
)

func TestMain(t *testing.M) {
	// r := NewRauzy(4)
	// p := map[int64][]int64{
	// 	0: {0, 1},
	// 	1: {0, 2},
	// 	2: {0, 3},
	// 	3: {0}}
	// r.SetSub(p)
	// r.Run(10)

	// r.Print()
	// r.Points("points.csv")

	r := NewRauzy(3)
	p := map[int64][]int64{
		0: {0, 1},
		1: {0, 2},
		2: {0}}
	r.SetSub(p)
	r.Run(25)

	r.Print()
	r.Png("rauzy.png")
}
