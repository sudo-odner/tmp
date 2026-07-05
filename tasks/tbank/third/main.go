package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)

	if !input.Scan() {
		return
	}
	n, _ := strconv.Atoi(input.Text())

	var b1 int64 = math.MinInt64
	var s1 int64 = int64(0)
	var b2 int64 = math.MinInt64
	var s2 int64 = int64(0)

	for i := 0; i < n; i++ {
		if !input.Scan() {
			break
		}
		price, _ := strconv.ParseInt(input.Text(), 10, 64)

		if -price > b1 {
			b1 = -price
		}
		if b1+price > s1 {
			s1 = b1 + price
		}
		if s1-price > b2 {
			b2 = s1 - price
		}
		if b2+price > s2 {
			s2 = b2 + price
		}
	}

	fmt.Println(s2)
}
