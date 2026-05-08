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
	sRune := []rune(s)
	if len(sRune) == 1 {
		return s
	}
	if len(sRune) == 2 {
		if sRune[0] == sRune[1] {
			return s
		}
		return string(sRune[0])
	}

	longest := sRune[:1]
	// Рассматриваем наибольшу подстроку у типов 'bab'
	for i := 0; i < len(sRune); i++ {
		for c := 1; i-c >= 0 && i+c < len(sRune) && sRune[i-c] == sRune[i+c]; c++ {
			if len(sRune[i-c:i+c+1]) > len(longest) {
				longest = sRune[i-c : i+c+1]
			}
		}
	}

	// Рассматриваем наибольшую под строкк к типов 'addf'
	for i := 1; i < len(sRune); i++ {
		for l, r := 1, 0; i-l >= 0 && i+r < len(sRune) && sRune[i-l] == sRune[i+r]; l, r = l+1, r+1 {
			if len(sRune[i-l:i+r+1]) > len(longest) {
				longest = sRune[i-l : i+r+1]
			}
		}
	}
	return string(longest)
}
