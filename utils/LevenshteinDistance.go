package utils

import (
	"unicode/utf8"
)

func LevenshteinDistance(a, b string) int {
	// Satisfy the basic condititons
	if a == b {
		return 0
	}

	if len(a) == 0 {
		return utf8.RuneCountInString(b)
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a)
	}

	s1 := []rune(a)
	s2 := []rune(b)

	x := make([]int, len(s1)+1)

	for i := 1; i < len(x); i++ {
		x[i] = int(i)
	}

	for i := 0; i < len(s2); i++ {
		prev := i + 1

		for j := 0; j < len(s1); j++ {

			current := x[j] //this is the matching case
			
			if s1[j] != s2[i] {
				current = min(
					x[j],
					prev,
					x[j+1],
				) + 1
			} 
		
			x[j] = prev
			prev = current
		}

		x[len(s1)] = prev
	}

	return x[len(s1)]
}

