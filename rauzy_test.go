package rauzy

import (
	"testing"
)

func TestMain(t *testing.M) {
	r := NewRauzy(4)
	p := map[int64][]int64{
		0: {0, 1},
		1: {0, 2},
		2: {0, 3},
		3: {0}}
	r.SetSub(p)
	r.Run(8)

	r.Print()
	r.Prj("points.csv")
	r.Png("rauzy.png")
}
