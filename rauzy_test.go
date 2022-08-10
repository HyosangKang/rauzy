package rauzy

import (
	"testing"

	"gonum.org/v1/plot"
)

func TestMain(t *testing.M) {
	r := NewRauzy([3][]int{{0, 1}, {0, 2}, {0}})
	bd := bdd(r.CPs)
	r.Run(2)
	r.Png(600, 600, bd, "img/rauzy_go.png")
	// r.Gif(600, 600, "img/rauzy.gif", 60)

	seq := []int{0, 1, 0, 2, 0, 1, 0}
	bseq := []int{0, 2, 1, 5, 3, 4, 0}
	f := NewFractal(
		r,
		seq,
		bseq,
		bd,
	)

	p := plot.New()
	setBdd(p, bd)
	f.Run(3)
	f.Points(p, false)
	p.Save(400, 400, "img/rauzy_bdp.png")

	// p := plot.New()
	// setBdd(p, bd)
	// f.Run(2)
	// f.Lines(p, true)
	// p.Save(400, 400, "img/rauzy_bdl_2.png")

}
