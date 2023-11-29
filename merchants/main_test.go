package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtractGoodsAndPrices(t *testing.T) {
	testCases := []struct {
		input          string
		expectedItems  []string
		expectedPrices []int
	}{
		{
			input:          "Bronze Cap 154-174 gil \nFaceguard 1334-1508 gil \nBronze Harness 235-266 gil",
			expectedItems:  []string{"Bronze Cap", "Faceguard", "Bronze Harness"},
			expectedPrices: []int{154, 174, 1334, 1508, 235, 266},
		},
		// Add more test cases for different scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// Replace newlines with spaces for more flexible comparison
			expectedItems := strings.Fields(strings.ReplaceAll(strings.TrimSpace(tc.input), "\n", " "))

			items, prices := extractGoodsAndPrices(tc.input)

			// Check if the extracted items and prices match the expected values
			if !reflect.DeepEqual(items, expectedItems) {
				t.Errorf("Expected items %v, but got %v", expectedItems, items)
			}

			if !reflect.DeepEqual(prices, tc.expectedPrices) {
				t.Errorf("Expected prices %v, but got %v", tc.expectedPrices, prices)
			}
		})
	}
}
