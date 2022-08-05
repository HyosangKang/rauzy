package rauzy

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type Fractal struct {
	R   *Rauzy
	Bdd [][2]float64
}

func NewFractal(r *Rauzy) *Fractal {
	return &Fractal{
		R:   r,
		Bdd: bdd(r.Pts),
	}
}

func (f *Fractal) Points(fn string) {
	mm := f.Bdd
	p := plot.New()
	xys := plotter.XYs{
		plotter.XY{X: mm[0][0], Y: mm[0][1]},
		plotter.XY{X: mm[1][0], Y: mm[1][1]},
	}
	s, _ := plotter.NewScatter(xys)
	s.GlyphStyle.Color = color.RGBA{255, 255, 255, 0}
	s.GlyphStyle.Radius = 0
	p.Add(s)

	l3 := len(f.R.Morph(3))
	l4 := len(f.R.Morph(4))
	_ = l4
	l5 := len(f.R.Morph(5))
	_ = l5
	l6 := len(f.R.Morph(6))
	_ = l6
	l9 := len(f.R.Morph(9))
	_ = l9
	for i := 0; i < l9; i++ {
		cp := f.R.Pts[i]
		xys = plotter.XYs{
			plotter.XY{X: cp.P[0], Y: cp.P[1]},
		}
		s, _ = plotter.NewScatter(xys)
		if i < l3 {
			s.GlyphStyle.Shape = draw.PyramidGlyph{}
		} else if i < l6 {
			s.GlyphStyle.Shape = draw.CircleGlyph{}
		} else {
			s.GlyphStyle.Shape = draw.RingGlyph{}
		}
		s.GlyphStyle.Color = f.R.Colors[cp.C]
		p.Add(s)

		lab := fmt.Sprintf("%d", i)
		label := plotter.XYLabels{
			XYs:    xys,
			Labels: []string{lab},
		}
		l, _ := plotter.NewLabels(label)
		p.Add(l)
	}

	ind := [][]int{}
	ind = append(ind, []int{0, 7, 13, 20, 24, 31, 37})
	ind = append(ind, []int{0, 44, 81, 125, 149, 193, 230})
	for _, in := range ind {
		for i := 1; i < len(in); i++ {
			fmt.Printf("%v ", f.R.Word[in[i-1]:in[i]])
		}
		fmt.Println()
	}

	v := [][2]float64{}
	for i := 0; i < 3; i++ {
		s := [3]float64{}
		s[i] = 1
		w := [2]float64{}
		for j := 0; j < 2; j++ {
			w[j] = dot(f.R.B[j], oprj(s, f.R.EV))
		}
		v = append(v, w)
	}

	v0 := f.R.Pts[0].P
	for i := 1; i < len(ind[0]); i++ {
		r := [2]float64{}
		for _, i := range f.R.Word[ind[0][i-1]:ind[0][i]] {
			for j := 0; j < 2; j++ {
				r[j] += v[i][j]
			}
		}
		xys = plotter.XYs{
			plotter.XY{X: v0[0], Y: v0[1]},
			plotter.XY{X: v0[0] + r[0], Y: v0[1] + r[1]},
		}
		l, _ := plotter.NewLine(xys)
		p.Add(l)
		for j := 0; j < 2; j++ {
			v0[j] += r[j]
		}
	}

	seq := []int{0, 1, 0, 2, 0, 1, 0}
	vv := [][2]float64{}
	for _, i := range seq {
		vv = append(vv, v[i])
	}
	for i := 0; i < 2; i++ {
		tvv := [][2]float64{}
		for _, j := range seq {
			t := [2]float64{}
			for k := 0; k < 2; k++ {
				if j < 2 {
					for q := 0; q < len(seq)-j; q++ {
						t[k] += vv[q][k]
					}
				} else {
					for q := 0; q < len(seq)-3; q++ {
						t[k] += vv[q][k]
					}
				}
			}
			tvv = append(tvv, t)
		}
		vv = tvv
	}

	v0 = f.R.Pts[15].P
	for i := 0; i < len(seq)-1; i++ {
		xys = plotter.XYs{
			plotter.XY{X: v0[0], Y: v0[1]},
			plotter.XY{X: v0[0] + vv[i][0], Y: v0[1] + vv[i][1]},
		}
		l, _ := plotter.NewLine(xys)
		p.Add(l)
		for j := 0; j < 2; j++ {
			v0[j] += vv[i][j]
		}
	}
	p.Save(400, 400, fn)
}
