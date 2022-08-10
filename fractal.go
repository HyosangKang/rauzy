package rauzy

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type Fractal struct {
	R         *Rauzy
	Seq, Bseq []int
	Bdd       [][2]float64
	FPs, BPs  map[int][]*FP
	VIs       map[int][][3]int
	Vs        map[int][][2]float64
	N         int // number of layers
}

type FP struct {
	XY  [2]float64 // projected vector
	V   [3]int     // word vector
	Ans *FP        // anscester
	Chd []*FP      // children
}

func NewFractal(r *Rauzy, seq, bseq []int, bdd [][2]float64) *Fractal {
	return &Fractal{
		R:    r,
		Seq:  seq,
		Bseq: bseq,
		Bdd:  bdd,
	}
}

func (f *Fractal) Init() {
	o := &FP{
		XY:  [2]float64{0, 0},
		Ans: nil,
		V:   [3]int{0, 0, 0},
	}
	f.FPs = make(map[int][]*FP)
	f.FPs[-1] = []*FP{o}
	f.BPs = make(map[int][]*FP)
	f.BPs[-1] = []*FP{o}

	v := [][2]float64{}
	vi := [][3]int{}
	for _, i := range f.Seq {
		s := [3]float64{}
		si := [3]int{}
		s[i] = 1
		si[i] = 1
		v = append(v, [2]float64{
			dot(f.R.B[0], oprj(s, f.R.EV)),
			dot(f.R.B[1], oprj(s, f.R.EV)),
		})
		vi = append(vi, si)
	}
	f.Vs = make(map[int][][2]float64)
	f.Vs[-1] = v

	f.VIs = make(map[int][][3]int)
	f.VIs[-1] = vi
}

func (f *Fractal) Run(n int) {
	f.Init()
	f.N = n
	f.R.Run(3 * n) // update rauzy
	for i := 0; i < n; i++ {
		f.Layer(i)
		f.Morph(i)
	}
}

// Adds the next layer of the rauzy points.
func (f *Fractal) Layer(n int) {
	f.FPs[n] = []*FP{}
	// updates the new points
	for _, fp := range f.FPs[n-1] {
		fps := []*FP{}
		p := fp.XY
		s := fp.V
		for i := 0; i < len(f.Seq); i++ {
			for j := 0; j < 2; j++ {
				p[j] += f.Vs[n-1][i][j]
			}
			for j := 0; j < 3; j++ {
				s[j] += f.VIs[n-1][i][j]
			}
			fps = append(fps, &FP{
				XY:  p,
				Ans: fp,
				V:   s,
			})
		}
		for _, i := range f.Bseq {
			fp.Chd = append(fp.Chd, fps[i])
			f.FPs[n] = addFP(f.FPs[n], fps[i])
		}
	}

	// updates the new boundary points
	f.BPs[n] = []*FP{}
	for _, bp := range f.BPs[n-1] {
		iid := -1
		for j, fp := range bp.Chd {
			if f.IsIn(fp, n) {
				iid = j
				break
			}
		}
		if iid == -1 {
			for _, fp := range bp.Chd {
				f.BPs[n] = addFP(f.BPs[n], fp)
			}
		} else {
			for j := 0; j < len(bp.Chd); j++ {
				k := (iid + j) % len(bp.Chd)
				if !f.IsIn(bp.Chd[k], n) {
					f.BPs[n] = addFP(f.BPs[n], bp.Chd[k])
				}
			}
		}
	}
}

func (f *Fractal) Morph(n int) {
	f.Vs[n] = [][2]float64{}
	f.VIs[n] = [][3]int{}
	for _, i := range f.Seq {
		t := [2]float64{}
		ti := [3]int{}
		m := i
		if i == 2 {
			m = 3
		}
		for k := 0; k < len(f.Seq)-m; k++ {
			for j := 0; j < 2; j++ {
				t[j] += f.Vs[n-1][k][j]
			}
			for j := 0; j < 3; j++ {
				ti[j] += f.VIs[n-1][k][j]
			}
		}
		f.Vs[n] = append(f.Vs[n], t)
		f.VIs[n] = append(f.VIs[n], ti)
	}
}

func addFP(fps []*FP, fp *FP) []*FP {
	for _, p := range fps {
		if equalVI(p.V, fp.V) {
			return fps
		}
	}
	return append(fps, fp)
}

func equalVI(p, q [3]int) bool {
	for i := 0; i < 3; i++ {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

func (fp *FP) Print() {
	fmt.Printf("XY: %v, (%d) Chd, (%v) Ans\n", fp.XY, len(fp.Chd), fp.Ans.XY)
}

func (fp *FP) AnsList() []*FP {
	if fp.Ans == nil {
		return []*FP{}
	}
	fpl := []*FP{fp.Ans}
	fpl = append(fpl, fp.Ans.AnsList()...)
	return fpl
}

func (f *Fractal) Boundary(n int) {
	f.BPs[n] = []*FP{}
	for _, fp := range f.FPs[n] {
		if !f.IsIn(fp, n) {
			f.BPs[n] = append(f.BPs[n], fp)
		}
	}
}

func (f *Fractal) IsIn(fp *FP, n int) bool {
	al := fp.AnsList()
	for _, afp := range al[1:] {
		for j := 1; j < len(afp.Chd); j++ {
			if tri(afp.XY, afp.Chd[j-1].XY, afp.Chd[j].XY, fp.XY) {
				return true
			}
		}
	}
	return false
}

// tri checks whether the vector p-c
// is in the triangle bounded by the vectors p0-c and p1-c
func tri(c, p0, p1, p [2]float64) bool {
	for i := 0; i < 2; i++ {
		p[i] -= c[i]
		p0[i] -= c[i]
		p1[i] -= c[i]
	}
	// if p0[0]*p1[1]-p0[1]*p1[0] <= 0 {
	// 	return false
	// }
	m := [2][2]float64{p0, p1}
	d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	a1 := (m[1][1]*p[0] - m[1][0]*p[1]) / d
	a2 := (-m[0][1]*p[0] + m[0][0]*p[1]) / d
	a := a1 + a2
	if 0 <= a1 && 0 <= a2 && a <= 1 {
		return true
	}
	return false
}

// commAns finds the closest common ansester of p0 and p1.
// func (f *Fractal) commAns(p0, p1 *FP) *FP {
// 	ap0 := p0.AnsList()
// 	ap1 := p1.AnsList()
// 	for _, p := range ap0 {
// 		for _, q := range ap1 {
// 			if p == q {
// 				return p
// 			}
// 		}
// 	}
// 	return nil
// }

func (f *Fractal) Points(p *plot.Plot, debug bool) {
	f.R.Points(p, 3*f.N, debug)

	xys := plotter.XYs{}
	for _, fp := range f.BPs[len(f.BPs)-2] {
		xys = append(xys, plotter.XY{
			X: fp.XY[0], Y: fp.XY[1]})
	}
	s, _ := plotter.NewScatter(xys)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(s)
}

func (f *Fractal) Lines(p *plot.Plot, debug bool) {
	f.R.Points(p, 3*f.N, debug)

	xys := plotter.XYs{}
	fps := f.BPs[len(f.BPs)-2]
	for _, fp := range fps {
		xys = append(xys, plotter.XY{
			X: fp.XY[0], Y: fp.XY[1]})
	}
	xys = append(xys, plotter.XY{
		X: fps[0].XY[0], Y: fps[0].XY[1]})
	l, _ := plotter.NewLine(xys)
	p.Add(l)
}

// func (f *Fractal) Png(fn string) {
// 	// Create new plot
// 	mm := f.Bdd
// 	p := plot.New()
// 	setBdd(p, mm)

// 	// Initialize sequence vectors
// 	seq := []int{0, 1, 0, 2, 0, 1, 0}
// 	bseq := []int{0, 2, 1, 5, 3, 4, 0}
// 	v := newSeqVec(seq, f.R.B, f.R.EV)

// 	bps := [][2]float64{{0, 0}}
// 	// For each loop 3*n
// 	n := 10
// 	for i := 0; i < n; i++ {
// 		// drawLine(p, bps)
// 		nbps := [][2]float64{}
// 		// For each boundary points
// 		for _, bp := range bps {
// 			// Create new boundary points
// 			nbp := [][2]float64{}
// 			for _, vv := range v {
// 				for k := 0; k < 2; k++ {
// 					bp[k] += vv[k]
// 				}
// 				nbp = append(nbp, bp)
// 			}
// 			if i == 0 {
// 				nbp = sortPts(nbp, bseq)
// 			} else {
// 				nbp = sortPts(nbp, bseq[:len(bseq)-1])
// 			}
// 			iid := -1
// 			// Search the last interior point
// 			for j, pt := range nbp {
// 				if interior(pt, bps) {
// 					iid = j
// 					break
// 				}
// 			}
// 			// Define new boundary
// 			if iid == -1 {
// 				// All points are exterior points
// 				// which happens only at i = 0
// 				nbps = append(nbps, nbp...)
// 			} else {
// 				for j := 0; j < len(seq)-1; j++ {
// 					k := (iid + j) % (len(seq) - 1)
// 					if !interior(nbp[k], bps) {
// 						nbps = append(nbps, nbp[k])
// 					}
// 				}
// 			}
// 		}
// 		rbps := [][2]float64{nbps[0]}
// 		for j := 1; j < len(nbps); j++ {
// 			if math.Abs(nbps[j-1][0]-nbps[j][0]) < 1e-10 {
// 				if math.Abs(nbps[j-1][1]-nbps[j][1]) < 1e-10 {
// 					fmt.Println("OK")
// 					continue
// 				}
// 			}
// 			rbps = append(rbps, nbps[j])
// 		}
// 		bps = rbps
// 		fmt.Println(len(bps))

// 		// Morph the sequence vectors
// 		tv := [][2]float64{}
// 		for _, j := range seq {
// 			t := [2]float64{}
// 			for k := 0; k < 2; k++ {
// 				if j < 2 {
// 					for q := 0; q < len(seq)-j; q++ {
// 						t[k] += v[q][k]
// 					}
// 				} else {
// 					for q := 0; q < len(seq)-3; q++ {
// 						t[k] += v[q][k]
// 					}
// 				}
// 			}
// 			tv = append(tv, t)
// 		}
// 		v = tv
// 	}
// 	drawLine(p, bps)
// 	wl := len(f.R.Morph(3 * n))
// 	for i := 0; i < wl; i++ {
// 		cp := f.R.CPs[i]
// 		xys := plotter.XYs{
// 			plotter.XY{X: cp.P[0], Y: cp.P[1]},
// 		}
// 		s, _ := plotter.NewScatter(xys)
// 		s.GlyphStyle.Shape = draw.CircleGlyph{}
// 		s.GlyphStyle.Color = f.R.Colors[cp.C]
// 		p.Add(s)

// 		// lab := fmt.Sprintf("%d", i)
// 		// label := plotter.XYLabels{
// 		// 	XYs:    xys,
// 		// 	Labels: []string{lab},
// 		// }
// 		// l, _ := plotter.NewLabels(label)
// 		// p.Add(l)
// 	}
// 	p.Save(1200, 1200, fn)
// }

func sortPts(pts [][2]float64, seq []int) [][2]float64 {
	np := [][2]float64{}
	for _, i := range seq {
		np = append(np, pts[i])
	}
	return np
}

func sortPts2(pts []*FP, seq []int) []*FP {
	np := []*FP{}
	for _, i := range seq {
		np = append(np, pts[i])
	}
	return np
}

// func newBdd(bps, nbps [][2]float64) [][2]float64 {
// 	nb := [][2]float64{}

// 	for _, bp := range nbps {
// 		if !interior(bp, bps) {
// 			nb = append(nb, bp)
// 		}
// 	}
// 	return nb
// }

func interior2(fp *FP, bps []*FP) bool {
	ret := false
	p := fp.XY
	for i := 1; i < len(bps); i++ {
		m := [2][2]float64{
			bps[i-1].XY, bps[i].XY}
		d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
		a1 := (m[1][1]*p[0] - m[1][0]*p[1]) / d
		a2 := (-m[0][1]*p[0] + m[0][0]*p[1]) / d
		a := a1 + a2
		if 0 <= a1 && 0 <= a2 && a <= 1 {
			ret = true
			goto NEXT
		}
	}
	return false
NEXT:

	return ret
}

func interior(p [2]float64, bps [][2]float64) bool {
	for i := 1; i < len(bps); i++ {
		m := [2][2]float64{bps[i-1], bps[i]}
		d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
		a1 := (m[1][1]*p[0] - m[1][0]*p[1]) / d
		a2 := (-m[0][1]*p[0] + m[0][0]*p[1]) / d
		a := a1 + a2
		if 0 <= a1 && 0 <= a2 && a <= 1 {
			return true
		}
	}
	return false
}

func newSeqVec(seq []int, b [][3]float64, ev [3]float64) (v [][2]float64) {
	for _, i := range seq {
		s := [3]float64{}
		s[i] = 1
		w := [2]float64{}
		for j := 0; j < 2; j++ {
			w[j] = dot(b[j], oprj(s, ev))
		}
		v = append(v, w)
	}
	return v
}

func drawLine2(p *plot.Plot, pts []*FP) {
	xys := plotter.XYs{}
	for _, pt := range pts {
		xys = append(xys, plotter.XY{
			X: pt.XY[0],
			Y: pt.XY[1],
		})
	}
	line, _ := plotter.NewLine(xys)
	p.Add(line)
}

func drawLine(p *plot.Plot, pts [][2]float64) {
	xys := plotter.XYs{}
	for _, pt := range pts {
		xys = append(xys, plotter.XY{
			X: pt[0],
			Y: pt[1],
		})
	}
	line, _ := plotter.NewLine(xys)
	p.Add(line)
}

func setBdd(p *plot.Plot, mm [][2]float64) {
	xys := plotter.XYs{
		plotter.XY{X: mm[0][0], Y: mm[0][1]},
		plotter.XY{X: mm[1][0], Y: mm[1][1]},
	}
	s, _ := plotter.NewScatter(xys)
	s.GlyphStyle.Color = color.RGBA{255, 255, 255, 0}
	s.GlyphStyle.Radius = 0
	p.Add(s)
}
