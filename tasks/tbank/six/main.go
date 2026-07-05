package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

type H []int64

func (h H) Len() int            { return len(h) }
func (h H) Less(i, j int) bool  { return h[i] < h[j] }
func (h H) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *H) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *H) Pop() interface{}   { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

const NEG = math.MinInt64 / 2

func main() {
	var a, b, n int
	if _, err := fmt.Scan(&a, &b, &n); err != nil {
		return
	}

	type C struct{ x, y int64 }
	cs := make([]C, n)
	for i := range cs {
		fmt.Scan(&cs[i].x, &cs[i].y)
	}

	sort.Slice(cs, func(i, j int) bool {
		return (cs[i].x - cs[i].y) > (cs[j].x - cs[j].y)
	})

	suffix := make([]int64, n+1)
	for i := range suffix {
		suffix[i] = NEG
	}
	if b == 0 {
		for i := range suffix {
			suffix[i] = 0
		}
	} else {
		h := &H{}
		var sum int64
		for i := n - 1; i >= 0; i-- {
			heap.Push(h, cs[i].y)
			sum += cs[i].y
			if h.Len() > b {
				sum -= heap.Pop(h).(int64)
			}
			if h.Len() == b {
				suffix[i] = sum
			}
		}
	}

	ans := int64(NEG)
	if a == 0 {
		for k := 0; k <= n-b; k++ {
			if suffix[k] != NEG && suffix[k] > ans {
				ans = suffix[k]
			}
		}
	} else {
		h := &H{}
		var prefixSum int64
		for i, c := range cs {
			heap.Push(h, c.x)
			prefixSum += c.x
			if h.Len() > a {
				prefixSum -= heap.Pop(h).(int64)
			}
			k := i + 1
			if h.Len() == a && k <= n-b && suffix[k] != NEG {
				if v := prefixSum + suffix[k]; v > ans {
					ans = v
				}
			}
		}
	}

	fmt.Println(ans)
}
