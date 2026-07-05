package main

import "fmt"

func main() {
	var i string
	if _, err := fmt.Scanln(&i); err != nil {
		return
	}

	n := len(i)

	lc := 0
	for lc < n && i[lc] == 'a' {
		lc++
	}
	if lc == n {
		fmt.Println("Yes")
		return
	}
	rc := 0
	for rc < n && i[n-1-rc] == 'a' {
		rc++
	}
	if lc > rc {
		fmt.Println("No")
		return
	}

	l := lc
	r := n - 1 - rc
	for l < r {
		if i[l] != i[r] {
			fmt.Println("No")
			return
		}
		l++
		r--
	}
	fmt.Println("Yes")
}
