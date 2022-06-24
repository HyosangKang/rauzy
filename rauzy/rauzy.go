package rauzy

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type rauzy struct {
	n     int
	seq   []int
	pisot func(int) []int
	vec   []float64
}

func NewRauzy(n int) *rauzy {
	return &rauzy{
		n:   n,
		seq: []int{0},
	}
}

func (r *rauzy) SetPisot(f func(int) []int) {
	for i := 0; i < r.n; i++ {
		if len(f(i)) == 0 {
			panic("Invalid pisot subsitution")
		}
	}
	r.pisot = f
}

func (r *rauzy) UpdateSeq(n int) {
	for i := 0; i < n; i++ {
		s := []int{}
		for _, ss := range r.seq {
			for _, sss := range r.pisot(ss) {
				s = append(s, sss)
			}
		}
		r.seq = s
	}
	v := make([]float64, r.n)
	for _, ss := range r.seq {
		v[ss] += 1
	}
	r.vec = v
}

func (r *rauzy) ClearSeq() {
	r.seq = []int{}
}

func (r *rauzy) SaveTxt(fn string) {
	r.print()
	fp, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	for _, s := range r.seq {
		fp.WriteString(fmt.Sprintf("%d ", s))
	}
}

func (r *rauzy) SavePng(fn string) {
	r.print()
	pvec := r.splitDots()
	basis := r.basis()
	pcoord := projCoord(pvec, basis)
	min, max := minmax(pcoord)
	trX, trY := trans(min, max, 600, 600)
	fp, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	pal := randomColor(r.n)
	pal = append(pal, color.White)
	width, height := 600, 600
	img := image.NewPaletted(image.Rect(0, 0, width, height), pal)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.White)
		}
	}
	for i, pc := range pcoord {
		for _, cc := range pc {
			img.Set(trX(cc[0]), trY(cc[1]), pal[i])
		}
	}
	png.Encode(fp, img)
}

func trans(min, max []float64, sx, sy float64) (func(float64) int, func(float64) int) {
	return func(x float64) int {
			return int(sx / (max[0] - min[0]) * (x - min[0]))
		}, func(y float64) int {
			return int(sy / (max[1] - min[1]) * (max[1] - y))
		}
}

func (r *rauzy) basis() [][]float64 {
	eb := [][]float64{}
	for i := 0; i < r.n; i++ {
		v := []float64{}
		for j := 0; j < r.n; j++ {
			if j == i {
				v = append(v, 1.0)
			} else {
				v = append(v, 0.0)
			}
		}
		eb = append(eb, v)
	}
	b := [][]float64{}
	for i := 0; i < r.n-1; i++ {
		if i == 0 {
			b = append(b, normalize(orthoProject(eb[i], r.vec)))
		} else {
			b = append(b, gramSchmidt(b, orthoProject(eb[i], r.vec)))
		}
	}
	return b
}

func minmax(pc [][][]float64) ([]float64, []float64) {
	n := len(pc[0][0])
	min := make([]float64, n)
	max := make([]float64, n)
	for i := 0; i < n; i++ {
		min[i] = math.MaxFloat64
		max[i] = -math.MaxFloat64
	}
	for _, cc := range pc {
		for _, c := range cc {
			for i := 0; i < n; i++ {
				if c[i] < min[i] {
					min[i] = c[i]
				}
				if c[i] > max[i] {
					max[i] = c[i]
				}
			}
		}
	}
	return min, max
}

func projCoord(pvec [][][]float64, basis [][]float64) [][][]float64 {
	pc := make([][][]float64, len(pvec))
	for i := 0; i < len(pvec); i++ {
		pc[i] = [][]float64{}
	}
	for i, vs := range pvec {
		for _, v := range vs {
			c := coord(v, basis)
			pc[i] = append(pc[i], c)
		}
	}
	return pc
}

func coord(v []float64, basis [][]float64) []float64 {
	c := []float64{}
	for _, b := range basis {
		c = append(c, inner(v, b))
	}
	return c
}

func gramSchmidt(basis [][]float64, v []float64) []float64 {
	w := []float64{}
	for _, vv := range v {
		w = append(w, vv)
	}
	for _, b := range basis {
		for i := 0; i < len(w); i++ {
			w[i] -= inner(w, b) / (norm(b) * norm(b)) * b[i]
		}
	}
	return normalize(w)
}

func (r *rauzy) splitDots() [][][]float64 {
	pvec := make([][][]float64, r.n)
	for i := 0; i < r.n; i++ {
		pvec[i] = [][]float64{}
	}
	v := make([]float64, r.n)
	for _, ss := range r.seq {
		v[ss] += 1.0
		pvec[ss] = append(pvec[ss], orthoProject(v, r.vec))
	}
	return pvec
}

func orthoProject(v, w []float64) []float64 {
	a := []float64{}
	for i, vv := range project(v, w) {
		a = append(a, v[i]-vv)
	}
	return a
}

func project(v, w []float64) []float64 {
	a := []float64{}
	k := inner(v, w) / (norm(w) * norm(w))
	for _, ww := range w {
		a = append(a, k*ww)
	}
	return a
}

func inner(v, w []float64) float64 {
	if len(v) != len(w) {
		panic("Invalid vector to project.")
	}
	a := 0.0
	for i := 0; i < len(v); i++ {
		a += v[i] * w[i]
	}
	return a
}

func normalize(v []float64) []float64 {
	w := []float64{}
	k := norm(v)
	for _, vv := range v {
		w = append(w, vv/k)
	}
	return w
}

func norm(v []float64) float64 {
	a := 0.0
	for _, vv := range v {
		a += float64(vv * vv)
	}
	return math.Sqrt(a)
}

func randomColor(n int) []color.Color {
	co := []color.Color{}
	for i := 0; i < n; i++ {
		r := uint8(rand.Intn(200))
		g := uint8(rand.Intn(200))
		b := uint8(rand.Intn(200))
		co = append(co, color.NRGBA{R: r, G: g, B: b, A: 255})
	}
	return co
}

func (r *rauzy) print() {
	fmt.Println()
	fmt.Println("::Pisot substitution::")
	for i := 0; i < r.n; i++ {
		fmt.Printf("%d -> %v\n", i, r.pisot(i))
	}
	fmt.Println("::Lenght of Fibonacci::")
	fmt.Printf("%d\n", len(r.seq))
	fmt.Println("::Pivot vector (normalized)::")
	fmt.Printf("%v (%v)\n\n", r.vec, normalize(r.vec))
}
