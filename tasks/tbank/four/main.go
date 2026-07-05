package main

import (
	"fmt"
)

const MOD = int64(1_000_000_007)

func powmod(base, exp int64) int64 {
	result := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			result = result * base % MOD
		}
		base = base * base % MOD
		exp /= 2
	}
	return result
}

func main() {
	var n, m int64
	fmt.Scan(&n, &m)

	mmod := m % MOD
	first := mmod * ((m - 1) % MOD) % MOD
	base := (mmod*mmod%MOD - 3*mmod%MOD + 3 + 2*MOD) % MOD

	fmt.Println(first * powmod(base, n-1) % MOD)
}
