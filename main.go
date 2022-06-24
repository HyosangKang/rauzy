package main

import "rauzy/rauzy"

func main() {
	r := rauzy.NewRauzy(3)
	p := func(n int) []int {
		switch n {
		case 0:
			return []int{0, 1}
		case 1:
			return []int{0, 2}
		case 2:
			return []int{0}
		}
		return []int{}
	}
	r.SetPisot(p)
	r.UpdateSeq(100)
	r.Save("rauzy.txt")
}
