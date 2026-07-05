package main

import (
	"fmt"
	"math"
)

func main() {
	var n int
	var x int64
	if _, err := fmt.Scan(&n, &x); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}
	const INF = int64(math.MaxInt64 / 2)

	m := [2]int64{0, INF}

	xc := x
	for i := 0; i < n-1; i++ {
		d := a[i+1] / a[i]
		r := xc % d
		xc /= d

		// fmt.Println(xc)
		mn := [2]int64{INF, INF}

		for c := int64(0); c <= 1; c++ {
			if m[c] == INF {
				continue
			}
			xi := r + c
			nc := int64(0)
			if xi >= d {
				xi -= d
				nc = 1
			}

			if mn[nc] > m[c]+xi {
				mn[nc] = m[c] + xi
			}
			if xi > 0 {
				if mn[nc+1] > m[c]+(d-xi) {
					mn[nc+1] = m[c] + (d - xi)
				}
			}
		}

		m = mn
	}

	answer := INF
	for c := int64(0); c <= 1; c++ {
		if m[c] == INF {
			continue
		}
		xi := xc + c
		if answer > m[c]+xi {
			answer = m[c] + xi
		}
	}

	fmt.Println(answer)
}
