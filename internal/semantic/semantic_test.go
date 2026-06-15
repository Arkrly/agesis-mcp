package semantic

import (
	"testing"
)

func FuzzDedupe(f *testing.F) {
	f.Add("a,b,c")
	f.Add("a,a,a")
	f.Add("")
	f.Add("a,b,a,c")

	f.Fuzz(func(t *testing.T, s string) {
		var input []string
		if s != "" {
			// Split by comma manually for fuzzing simplicity
			start := 0
			for i := 0; i < len(s); i++ {
				if s[i] == ',' {
					input = append(input, s[start:i])
					start = i + 1
				}
			}
			input = append(input, s[start:])
		}

		result := dedupe(input)

		if len(input) == 0 && len(result) != 0 {
			t.Errorf("dedupe of empty slice should be empty, got %d items", len(result))
		}

		// Check for duplicates in the result
		seen := make(map[string]bool)
		for _, item := range result {
			if seen[item] {
				t.Errorf("duplicate item %q found in result", item)
			}
			seen[item] = true
		}

		// Check if all unique items from input are in result
		inputMap := make(map[string]bool)
		for _, item := range input {
			inputMap[item] = true
		}
		if len(result) != len(inputMap) {
			t.Errorf("result length %d does not match number of unique input items %d", len(result), len(inputMap))
		}
	})
}
