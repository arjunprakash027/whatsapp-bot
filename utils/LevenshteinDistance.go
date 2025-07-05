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

		// this is the inserstion cost at the start
		prev := i + 1

		for j := 0; j < len(s1); j++ {

			// this is the current cost to edit S2[j] to S2[i]
			current := x[j] //this is the matching case
			
			if s1[j] != s2[i] {
				current = min(
					x[j], // this is the substitution cost
					prev, // this as said before is the insertion cost
					x[j+1], //This is the deletion cost
				) + 1
			} 
		
			x[j] = prev
			prev = current
		}

		x[len(s1)] = prev
	}

	return x[len(s1)]
}

