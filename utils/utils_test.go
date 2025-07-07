package utils

import (
	"fmt"
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {

	got := NormalizedLevenshteinDistance("cat is where", "cut is here")
	fmt.Println("LevenshteinDistance = ", got)

	string_normalize := NormalizeText("CARLY \n does \r not    know why       🤔 ")

	fmt.Println("Normalized text = ", string_normalize)
	// t.Run(
	// 	"LevenshteinDistance",
	// 	func(t *testing.T) {
	// 		got := NormalizedLevenshteinDistance("cat","cut")

	// 		if got != 1 {
	// 			t.Errorf("LevenshteinDistance = %d; want 1", got)
	// 		}
	// 	},
	// )
}
