package rauzy

import (
	"testing"
)

func TestMain(t *testing.M) {
	r := NewRauzy(3)
	p := map[int64][]int64{
		0: {0, 1},
		1: {0, 2},
		2: {0}}
	r.SetPisot(p)
	r.UpdateSeq(25)
	r.SavePng("rauzy.png")
}
