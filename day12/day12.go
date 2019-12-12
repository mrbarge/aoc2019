package main

import (
	"fmt"
	"github.com/ernestosuarez/itertools"
	"math"
)

type Position struct {
	x float64
	y float64
	z float64
}

type Velocity struct {
	x float64
	y float64
	z float64
}

func step(p []*Position, v []Velocity) {
	// get all permutations of moons
	r := 2
	iterable := []int{0,1,2,3}

	diffs := make([]Position, 0)
	for i := 0; i < len(p); i++ {
		diffs = append(diffs, Position{0,0,0})
	}

	for v := range itertools.PermutationsInt(iterable, r) {
		m1 := p[v[0]]
		m2 := p[v[1]]

		if m1.x < m2.x {
			diffs[v[0]].x += 1
			diffs[v[1]].x -= 1
		} else if m1.x > m2.x {
			diffs[v[0]].x -= 1
			diffs[v[1]].x += 1
		}
		if m1.y < m2.y {
			diffs[v[0]].y += 1
			diffs[v[1]].y -= 1
		} else if m1.y > m2.y {
			diffs[v[0]].y -= 1
			diffs[v[1]].y += 1
		}
		if m1.z < m2.z {
			diffs[v[0]].z += 1
			diffs[v[1]].z -= 1
		} else if m1.z > m2.z {
			diffs[v[0]].z -= 1
			diffs[v[1]].z += 1
		}
	}

	// move the moons
	for i, _ := range diffs {
		fmt.Printf("Applying %d %d %d\n", diffs[i].x, diffs[i].y, diffs[i].z)
		p[i].x += diffs[i].x
		p[i].y += diffs[i].y
		p[i].z += diffs[i].z
	}

	// apply velocity
	for i, _ := range v {
		p[i].x += v[i].x
		p[i].y += v[i].y
		p[i].z += v[i].z
	}
}

func energy(p []*Position, v []Velocity) (r []float64) {
	r = make([]float64, len(p))
	for i, _ := range p {
		r[i] = (math.Abs(p[i].x) + math.Abs(p[i].y) + math.Abs(p[i].z)) * (math.Abs(v[i].x) + math.Abs(v[i].y) + math.Abs(v[i].z))
	}
	return r
}

func test() ([]*Position, []Velocity) {
	p := make([]*Position, 0)
	v := make([]Velocity, 0)
	p = append(p, &Position{-1, 0, 2})
	p = append(p, &Position{2, -10, 7})
	p = append(p, &Position{4,-8,8})
	p = append(p, &Position{3,5,-1})
	v = append(v, Velocity{0,0,0} )
	v = append(v, Velocity{0,0,0})
	v = append(v, Velocity{0,0,0})
	v = append(v, Velocity{0,0,0})
	return p, v
}

func input() ([]*Position, []Velocity) {
	p := make([]*Position, 0)
	v := make([]Velocity, 0)
	p = append(p, &Position{-7, -8, 9})
	p = append(p, &Position{-12, -3, -4})
	p = append(p, &Position{6,-17,-9})
	p = append(p, &Position{4,-10,-6})
	v = append(v, Velocity{0,0,0} )
	v = append(v, Velocity{0,0,0})
	v = append(v, Velocity{0,0,0})
	v = append(v, Velocity{0,0,0})
	return p, v
}

func main() {

	p, v := test()
	for i := 0; i < 1000; i++ {
		step(p, v)
	}
	fmt.Println(energy(p, v))
}
