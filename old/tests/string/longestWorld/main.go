package main

import "fmt"

// Напиши функцию на Go, которая принимает строку и возвращает самое длинное слово в этой строке. Слова разделены пробелами.
func longestWorld(s string) string {
	var answer string
	wordTmp := make([]rune, 0, len(s))

	for _, r := range s {
		if r == ' ' {
			if len(wordTmp) > len(answer) {
				answer = string(wordTmp)
			}
			wordTmp = wordTmp[:0]
		} else {
			wordTmp = append(wordTmp, r)
		}
	}

	if len(wordTmp) > len(answer) {
		answer = string(wordTmp)
	}

	return answer
}

func main() {
	fmt.Println(longestWorld("test testtest testest"))
}
