package ffxi

import (
	"reflect"
	"testing"
)

func TestExtractItemsFromIngredients(t *testing.T) {
	tests := []struct {
		ingredientsString string
		expectedItems     []Item
	}{
		{
			"HQ1: Antidote x6\nHQ2: Antidote x9\nHQ3: Antidote x12",
			[]Item{
				{Name: "Antidote", Count: 6},
				{Name: "Antidote", Count: 9},
				{Name: "Antidote", Count: 12},
			},
		},
		{
			"HQ1: Maple Shield +1\n",
			[]Item{
				{Name: "Maple Shield +1", Count: 1},
			},
		},
		{
			"HQ1: Angler's Hose\n",
			[]Item{
				{Name: "Angler's Hose", Count: 1},
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := extractItemsFromIngredients(test.ingredientsString)
		if !reflect.DeepEqual(result, test.expectedItems) {
			t.Errorf("For ingredients string %s, expected %v, but got %v", test.ingredientsString, test.expectedItems, result)
		}
	}
}
