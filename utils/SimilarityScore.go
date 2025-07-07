package utils

import (
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

func NormalizedLevenshteinDistance(a, b string) float64 {
	// Satisfy the basic condititons
	if a == b {
		return 0
	}

	if len(a) == 0 {
		return float64(utf8.RuneCountInString(b))
	}

	if len(b) == 0 {
		return float64(utf8.RuneCountInString(a))
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
					x[j],   // this is the substitution cost
					prev,   // this as said before is the insertion cost
					x[j+1], //This is the deletion cost
				) + 1
			}

			x[j] = prev
			prev = current
		}

		x[len(s1)] = prev
	}

	maxLength := max(len(s1), len(s2))
	log.Println("Max length", maxLength)
	log.Println("Distnace = ", x[len(s1)])
	return float64(x[len(s1)]) / float64(maxLength)
}

func NormalizeText(text string) string {

	text = strings.ToLower(text)

	// Remove emojis
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F900}-\x{1F9FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)
	text = emojiRegex.ReplaceAllString(text, "")

	//put everything under single line
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	// remove the extra spaces, the loop is required because the replacement happens only once and it must be iterated again
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	text = strings.TrimSpace(text)

	return text
}
