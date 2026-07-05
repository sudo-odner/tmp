package main

import "fmt"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}

	nn := make([][]int, n)
	for i := 0; i < n; i++ {
		nn[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		if s[i] == 'T' || s[i] == 'O' || s[i] == 'I' {
			nn[i][i] = 0
		} else {
			nn[i][i] = 1
		}
	}

	simb := []byte{'T', 'O', 'I'}
	for lh := 2; lh <= n; lh++ {
		for i := 0; i <= n-lh; i++ {
			j := i + lh - 1

			res := min(nn[i+1][j]+1, nn[i][j-1]+1)

			for _, c := range simb {
				ci := 0
				if s[i] != c {
					ci = 1
				}
				cj := 0
				if s[j] != c {
					cj = 1
				}

				ic := 0
				if i+1 <= j-1 {
					ic = nn[i+1][j-1]
				}

				res = min(res, ci+cj+ic)
			}

			nn[i][j] = res
		}
	}

	fmt.Println(nn[0][n-1])
}
