package utils

import (
	"testing"
	"fmt"
)

func TestLevenshteinDistance(t *testing.T) {
	
	got := LevenshteinDistance("cat is where","cut is here")
	fmt.Println("LevenshteinDistance = ", got)

	t.Run(
		"LevenshteinDistance",
		func(t *testing.T) {
			got := LevenshteinDistance("cat","cut")

			if got != 1 {
				t.Errorf("LevenshteinDistance = %d; want 1", got)
			}
		},
	)
}