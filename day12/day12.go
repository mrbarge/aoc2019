package main

import (
	"fmt"
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

func getPerms() [][]int {
	return [][]int{
		[]int{0,1}, []int{0,2},[]int{0,3},[]int{1,2},[]int{1,3},[]int{2,3},
	}
}

func step(p []*Position, v []*Velocity) {
	//fmt.Printf("Beginning 0 with %f %f %f and v %f %f %f\n", p[0].x, p[0].y, p[0].z, v[0].x, v[0].y, v[0].z)
	//fmt.Printf("Beginning 1 with %f %f %f and v %f %f %f\n", p[1].x, p[1].y, p[1].z, v[1].x, v[1].y, v[1].z)
	//fmt.Printf("Beginning 2 with %f %f %f and v %f %f %f\n", p[2].x, p[2].y, p[2].z, v[2].x, v[2].y, v[2].z)
	//fmt.Printf("Beginning 3 with %f %f %f and v %f %f %f\n", p[3].x, p[3].y, p[3].z, v[3].x, v[3].y, v[3].z)

	diffs := make([]Position, 0)
	for i := 0; i < len(p); i++ {
		diffs = append(diffs, Position{0,0,0})
	}

	for _, val := range getPerms() {
		m1 := p[val[0]]
		m2 := p[val[1]]

		//fmt.Printf("Looking at %d-%d\n", val[0], val[1])
		if m1.x < m2.x {
			diffs[val[0]].x += 1
			diffs[val[1]].x -= 1
		} else if m1.x > m2.x {
			diffs[val[0]].x -= 1
			diffs[val[1]].x += 1
		}
		if m1.y < m2.y {
			diffs[val[0]].y += 1
			diffs[val[1]].y -= 1
		} else if m1.y > m2.y {
			diffs[val[0]].y -= 1
			diffs[val[1]].y += 1
		}
		if m1.z < m2.z {
			diffs[val[0]].z += 1
			diffs[val[1]].z -= 1
		} else if m1.z > m2.z {
			diffs[val[0]].z -= 1
			diffs[val[1]].z += 1
		}
	}

	// move the moons
	//fmt.Printf("Applying %f %f %f\n", diffs[0].x, diffs[0].y, diffs[0].z)
	for i, _ := range diffs {
		v[i].x += diffs[i].x
		v[i].y += diffs[i].y
		v[i].z += diffs[i].z
	}

	// apply velocity
	for i, _ := range v {
		p[i].x += v[i].x
		p[i].y += v[i].y
		p[i].z += v[i].z
	}
	//fmt.Printf("Ending with %f %f %f\n", p[0].x, p[0].y, p[0].z)
}

func energy(p []*Position, v []*Velocity) (r []float64) {
	r = make([]float64, len(p))
	for i, _ := range p {
		r[i] = (math.Abs(p[i].x) + math.Abs(p[i].y) + math.Abs(p[i].z)) * (math.Abs(v[i].x) + math.Abs(v[i].y) + math.Abs(v[i].z))
	}
	return r
}

func allDone(b []bool) bool {
	for _, v := range b {
		if !v {
			return false
		}
	}
	return true
}

func findStepsUntilRepeat(p []*Position, v []*Velocity) (c []int64) {
	done := make([]bool, len(p))
	stepCount := make([]int64, len(p))
	starting := make([]Position, len(p))
	steps := int64(0)
	for _, val := range p {
		//done = append(done, false)
		starting = append(starting, Position{val.x, val.y, val.z})
	}

	for !allDone(done) {
		step(p, v)
		for i, _ := range p {
			if *p[i] == starting[i] && !done[i] {
				done[i] = true
				stepCount[i] = steps
				fmt.Printf("Done %d\n", i)
			}
		}
		steps += 1
	}

	return stepCount
}

func test() ([]*Position, []*Velocity) {
	p := make([]*Position, 0)
	v := make([]*Velocity, 0)
	p = append(p, &Position{-1, 0, 2})
	p = append(p, &Position{2, -10, -7})
	p = append(p, &Position{4,-8,8})
	p = append(p, &Position{3,5,-1})
	v = append(v, &Velocity{0,0,0} )
	v = append(v, &Velocity{0,0,0})
	v = append(v, &Velocity{0,0,0})
	v = append(v, &Velocity{0,0,0})
	return p, v
}

func input() ([]*Position, []*Velocity) {
	p := make([]*Position, 0)
	v := make([]*Velocity, 0)
	p = append(p, &Position{-7, -8, 9})
	p = append(p, &Position{-12, -3, -4})
	p = append(p, &Position{6,-17,-9})
	p = append(p, &Position{4,-10,-6})
	v = append(v, &Velocity{0,0,0} )
	v = append(v, &Velocity{0,0,0})
	v = append(v, &Velocity{0,0,0})
	v = append(v, &Velocity{0,0,0})
	return p, v
}

func main() {

	p, v := input()
	for i := 0; i < 1000; i++ {
		step(p, v)
	}

	fmt.Println(energy(p, v))

	p, v = input()
	stepCount := findStepsUntilRepeat(p, v)
	fmt.Println(stepCount)
}
