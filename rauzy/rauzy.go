package rauzy

import (
	"fmt"
	"os"
)

type rauzy struct {
	n     int
	seq   []int
	pisot func(int) []int
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
}

func (r *rauzy) ClearSeq() {
	r.seq = []int{}
}

func (r *rauzy) Save(fn string) {
	fp, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	for _, s := range r.seq {
		fp.WriteString(fmt.Sprintf("%d ", s))
	}
}
