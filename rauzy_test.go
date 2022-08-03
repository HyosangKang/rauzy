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
	// r.Points("img/points.csv")

	r := NewRauzy(20, [3][]int{{0, 1}, {0, 2}, {0}})
	r.Png(600, 600, nil, "img/rauzy_go_1_new.png")
	r.Png(600, 600,
		[][2]float64{
			{-.05, -.05},
			{.05, .05}},
		"img/rauzy_go_2.png")
	r.Run(2)
	r.Png(600, 600,
		[][2]float64{
			{-.005, -.005},
			{.005, .005}},
		"img/rauzy_go_3.png")
	r.Gif(600, 600, "img/rauzy_mov.gif", 60)
}
