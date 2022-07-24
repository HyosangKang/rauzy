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

type Rauzy struct {
	Dim    int
	Seq    []int64
	Sub    map[int64][]int64
	Vec    []float64
	Colors []color.Color
	Basis  [][]float64
	Pvec   map[int64][][]float64
}

func NewRauzy(n int) *Rauzy {
	co := []color.Color{}
	for i := 0; i < n; i++ {
		r := uint8(rand.Intn(200))
		g := uint8(rand.Intn(200))
		b := uint8(rand.Intn(200))
		co = append(co, color.NRGBA{R: r, G: g, B: b, A: 255})
	}
	co = append(co, color.White)

	return &Rauzy{
		Dim:    n,
		Seq:    []int64{0},
		Pvec:   make(map[int64][][]float64),
		Colors: co,
	}
}

func (r *Rauzy) SetSub(s map[int64][]int64) {
	for i := 0; i < r.Dim; i++ {
		if len(s[int64(i)]) == 0 {
			panic("Invalid pisot subsitution")
		}
	}
	r.Sub = s
}

func (r *Rauzy) Run(n int) {
	// Find sequence
	for i := 0; i < n; i++ {
		seq := []int64{}
		for _, s := range r.Seq {
			seq = append(seq, r.Sub[s]...)
		}
		r.Seq = seq
	}

	// Find eigenvector
	v := make([]float64, r.Dim)
	for _, ss := range r.Seq {
		v[ss] += 1
	}
	r.Vec = nrmz(v)

	// Find orthonormal basis perp to ev
	eb := [][]float64{}
	for i := 0; i < r.Dim-1; i++ {
		v := []float64{}
		for j := 0; j < r.Dim; j++ {
			if j == i {
				v = append(v, 1)
			} else {
				v = append(v, 0)
			}
		}
		eb = append(eb, v)
	}
	r.Basis = append(r.Basis, nrmz(oprj(eb[0], r.Vec)))
	for i := 1; i < r.Dim-1; i++ {
		r.Basis = append(r.Basis, gramSchmidt(r.Basis, oprj(eb[i], r.Vec)))
	}

	fmt.Println(r.Basis)
	for _, v := range r.Basis {
		for _, w := range r.Basis {
			fmt.Println(dot(v, w))
		}
	}

	// Find projections
	vg := make(map[int64][][]float64)
	v = make([]float64, r.Dim)
	for _, ss := range r.Seq {
		v[ss] += 1.0
		vg[ss] = append(vg[ss], oprj(v, r.Vec))
	}

	for i, vs := range vg {
		for _, v := range vs {
			c := coord(v, r.Basis)
			r.Pvec[i] = append(r.Pvec[i], c)
		}
	}
}

func (r *Rauzy) Points(fn string) {
	fp, _ := os.Create(fn)
	defer fp.Close()

	// add key row
	ax := []string{"x", "y", "z"}
	t := ""
	for i := 0; i < r.Dim; i++ {
		for j := 0; j < r.Dim-1; j++ {
			t += ax[j] + fmt.Sprintf("%d,", i)
		}
	}
	fp.WriteString(t)

	count := 0
	for {
		t = "\n"
		empty := true
		for i := 0; i < r.Dim; i++ {
			if len(r.Pvec[int64(i)]) > count {
				empty = false
				for _, v := range r.Pvec[int64(i)][count] {
					t += fmt.Sprintf("%f,", v)
				}
			} else {
				for j := 0; j < r.Dim-1; j++ {
					t += ","
				}
			}
		}
		if empty {
			break
		}
		fp.WriteString(t)
		count++
	}
}

func (r *Rauzy) Png(fn string) {
	min, max := bds(r.Pvec)
	trX, trY := trans(min, max, 600, 600)

	fp, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	width, height := 600, 600
	img := image.NewPaletted(image.Rect(0, 0, width, height), r.Colors)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.White)
		}
	}
	for i, cs := range r.Pvec {
		for _, c := range cs {
			img.Set(trX(c[0]), trY(c[1]), r.Colors[i])
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

func bds(pc map[int64][][]float64) ([]float64, []float64) {
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

func coord(v []float64, basis [][]float64) []float64 {
	c := []float64{}
	for _, b := range basis {
		c = append(c, dot(v, b))
	}
	return c
}

func gramSchmidt(basis [][]float64, v []float64) []float64 {
	w := make([]float64, len(v))
	copy(w, v)
	for _, b := range basis {
		ww := make([]float64, len(v))
		copy(ww, w)
		for i := 0; i < len(w); i++ {
			w[i] -= dot(ww, b) / (norm(b) * norm(b)) * b[i]
		}
	}
	return nrmz(w)
}

func oprj(v, w []float64) []float64 {
	a := []float64{}
	for i, vv := range prj(v, w) {
		a = append(a, v[i]-vv)
	}
	return a
}

func prj(v, w []float64) []float64 {
	a := []float64{}
	k := dot(v, w) / (norm(w) * norm(w))
	for _, ww := range w {
		a = append(a, k*ww)
	}
	return a
}

func dot(v, w []float64) float64 {
	if len(v) != len(w) {
		panic("Invalid vector to project.")
	}
	a := 0.0
	for i := 0; i < len(v); i++ {
		a += v[i] * w[i]
	}
	return a
}

func nrmz(v []float64) []float64 {
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

func (r *Rauzy) Print() {
	fmt.Println()
	fmt.Println("::Pisot substitution::")
	for i := 0; i < r.Dim; i++ {
		fmt.Printf("%d -> %v\n", i, r.Sub[int64(i)])
	}

	fmt.Println()
	fmt.Println("::Lenght of Fibonacci::")
	fmt.Printf("%d\n", len(r.Seq))

	fmt.Println()
	fmt.Println("::Pivot vector (normalized)::")
	fmt.Printf("%v\n\n", r.Vec)
}
