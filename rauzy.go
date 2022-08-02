package rauzy

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
)

type CP struct {
	C int        // color
	P [2]float64 // coordinate
}

type Rauzy struct {
	N      int           // max number of sequence
	Word   []int         // word sequence
	Sub    [3][]int      // substitution rule
	EV     [3]float64    // diverging direction
	B      [][3]float64  // basis of contracting plane
	Colors []color.Color // color set
	Pts    []CP          // Projected points
}

func NewRauzy(n int, s [3][]int) *Rauzy {
	for i := 0; i < 3; i++ {
		if len(s[i]) == 0 {
			panic("Invalid pisot subsitution")
		}
	}
	co := []color.Color{
		color.RGBA{245, 78, 66, 255},
		color.RGBA{90, 245, 66, 255},
		color.RGBA{38, 55, 237, 255},
		color.White,
		color.Black,
	}

	return &Rauzy{
		N:      n,
		Word:   []int{0},
		Colors: co,
		Sub:    s,
	}
}

func (r *Rauzy) Morph() {
	for i := 0; i < r.N; i++ {
		w := []int{}
		for _, a := range r.Word {
			w = append(w, r.Sub[a]...)
		}
		r.Word = w
	}
}

func (r *Rauzy) Eigenvector() {
	v := [3]float64{}
	for _, ss := range r.Word {
		v[ss] += 1
	}
	r.EV = nrmz(v)
}

func (r *Rauzy) Basis() {
	e1 := nrmz(oprj([3]float64{1, 0, 0}, r.EV))
	e2 := oprj([3]float64{0, 1, 0}, r.EV)
	k := dot(e1, e2)
	for i := 0; i < 3; i++ {
		e2[i] -= k * e1[i]
	}
	e2 = nrmz(e2)
	r.B = [][3]float64{e1, e2}
}

func (r *Rauzy) Run() {
	r.Morph()
	r.Eigenvector()
	r.Basis()
	r.Project()
}

func (r *Rauzy) Project() {
	v := [3]float64{0, 0, 0}
	for _, a := range r.Word {
		v[a] += 1.0
		o := oprj(v, r.EV)
		c := [2]float64{dot(o, r.B[0]), dot(o, r.B[1])}
		r.Pts = append(r.Pts, CP{a, c})
	}
}

func (r *Rauzy) Png(w, h int, fn string) {
	min, max := bds(r.Pts)
	trX, trY := trans(min, max, float64(w), float64(h))

	fp, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	img := image.NewPaletted(image.Rect(0, 0, w, h), r.Colors)
	clear(img)
	for _, cp := range r.Pts {
		img.Set(trX(cp.P[0]), trY(cp.P[1]), r.Colors[cp.C])
	}
	png.Encode(fp, img)
}

func (r *Rauzy) Gif(w, h int, fn string, sec int) {
	drawP := func(img *image.Paletted, p [2]int, c color.Color) {
		sz := 3
		for i := -sz; i <= sz; i++ {
			for j := -sz; j <= sz; j++ {
				img.Set(p[0]+i, p[1]+j, c)
			}
		}
	}
	drawL := func(img *image.Paletted, p, q [2]int) {
		n := p[0] - q[0]
		if n < 0 {
			n = -n
		}
		m := p[1] - q[1]
		if m < 0 {
			m = -n
		}
		if n < m {
			n = m
		}
		for i := 0; i < n; i++ {
			x := p[0] + int(float64(i*(q[0]-p[0]))/float64(n))
			y := p[1] + int(float64(i*(q[1]-p[1]))/float64(n))
			img.Set(x, y, color.Black)
		}
	}
	min, max := bds(r.Pts)
	trX, trY := trans(min, max, float64(w), float64(h))
	fp, _ := os.Create(fn)
	defer fp.Close()

	g := &gif.GIF{
		LoopCount: -1,
	}

	img := image.NewPaletted(image.Rect(0, 0, w, h), r.Colors)
	clear(img)
	for i := 0; i < 60*sec; i++ {
		if i > len(r.Pts)-1 {
			break
		}
		p := [2]int{trX(r.Pts[i].P[0]), trY(r.Pts[i].P[1])}
		c := r.Colors[r.Pts[i].C]
		drawP(img, p, c)
		tim := image.NewPaletted(image.Rect(0, 0, w, h), r.Colors)
		copy(tim.Pix, img.Pix)
		q := [2]int{trX(r.Pts[i+1].P[0]), trY(r.Pts[i+1].P[1])}
		drawL(tim, p, q)
		g.Image = append(g.Image, tim)
		g.Delay = append(g.Delay, 1)
	}
	gif.EncodeAll(fp, g)
}

func clear(img *image.Paletted) {
	w, h := img.Rect.Max.X, img.Rect.Max.Y
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.Set(i, j, color.White)
		}
	}
}

func trans(min, max []float64, sx, sy float64) (func(float64) int, func(float64) int) {
	dx, dy := max[0]-min[0], max[1]-min[1]
	if dy/dx > sy/sx {
		w := sy * dx / dy
		m := (sx - w) / 2
		return func(x float64) int {
				return int(w/dx*(x-min[0]) + m)
			}, func(y float64) int {
				return int(sy / dy * (max[1] - y))
			}
	} else {
		h := sx * dy / dx
		m := (sy - h) / 2
		return func(x float64) int {
				return int(sx / dx * (x - min[0]))
			}, func(y float64) int {
				return int(h/dy*(max[1]-y) + m)
			}
	}
}

func bds(pts []CP) ([]float64, []float64) {
	n := len(pts[0].P)
	min := make([]float64, n)
	max := make([]float64, n)
	for i := 0; i < n; i++ {
		min[i] = math.MaxFloat64
		max[i] = -math.MaxFloat64
	}
	for _, p := range pts {
		pp := p.P
		for i := 0; i < n; i++ {
			if pp[i] < min[i] {
				min[i] = pp[i]
			}
			if pp[i] > max[i] {
				max[i] = pp[i]
			}
		}
	}
	return min, max
}

func oprj(v, w [3]float64) [3]float64 {
	a := [3]float64{}
	for i, vv := range prj(v, w) {
		a[i] = v[i] - vv
	}
	return a
}

func prj(v, w [3]float64) []float64 {
	a := []float64{}
	k := dot(v, w) / (norm(w) * norm(w))
	for _, ww := range w {
		a = append(a, k*ww)
	}
	return a
}

func dot(v, w [3]float64) float64 {
	a := 0.0
	for i := 0; i < 3; i++ {
		a += v[i] * w[i]
	}
	return a
}

func nrmz(v [3]float64) [3]float64 {
	w := [3]float64{}
	k := norm(v)
	for i, vv := range v {
		w[i] = vv / k
	}
	return w
}

func norm(v [3]float64) float64 {
	a := 0.0
	for _, vv := range v {
		a += float64(vv * vv)
	}
	return math.Sqrt(a)
}

func (r *Rauzy) Print() {
	fmt.Println()
	fmt.Println("::Pisot substitution::")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d -> %v\n", i, r.Sub[i])
	}

	fmt.Println()
	fmt.Println("::Lenght of Fibonacci::")
	fmt.Printf("%d\n", len(r.Word))

	fmt.Println()
	fmt.Println("::Pivot vector (normalized)::")
	fmt.Printf("%v\n\n", r.EV)
}
