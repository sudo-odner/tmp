package main

// Given a string s, return the longest in s.
//
//
// Example 1:
//
// Input: s = "babad"
// Output: "bab"
// Explanation: "aba" is also a valid answer.
//
// Example 2:
//
// Input: s = "cbbd"
// Output: "bb"
//
//
// Constraints:
//     1 <= s.length <= 1000
//     s consist of only digits and English letters.

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}
	end, start := 0, 0
	for i := 0; i < len(s); i++ {
		l1, r1 := expand(s, i, i)
		l2, r2 := expand(s, i, i+1)
		if r1-l1 > end-start {
			start, end = l1, r1
		}
		if r2-l2 > end-start {
			start, end = l2, r2
		}
	}
	return s[start : end+1]
}

func expand(s string, l, r int) (int, int) {
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}
	return l + 1, r - 1
}
