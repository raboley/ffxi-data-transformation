package recipe

import (
	"reflect"
	"testing"
)

// TestExtractCraftType is a unit test for extractCraftType.
func TestExtractCraftType(t *testing.T) {
	tests := []struct {
		text     string
		expected string
	}{
		{"Guild Recipes: Woodworking", "Woodworking"},
		{"Guild Recipes: Alchemy", "Alchemy"},
		{"Guild Recipes: Smithing", "Smithing"},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := extractCraftType(test.text)
		if result != test.expected {
			t.Errorf("For text %s, expected %s, but got %s", test.text, test.expected, result)
		}
	}
}

// TestOtherRequirementSkillLevels is a unit test for otherRequirementSkillLevels.
func TestOtherRequirementSkillLevels(t *testing.T) {
	tests := []struct {
		requirements string
		expected     map[string]int
	}{
		{"Apprentice\nAlchemy(49)\n", map[string]int{"Alchemy": 49}},
		{"Apprentice\nAlchemy(30)\nSmithing(20)\n", map[string]int{"Alchemy": 30, "Smithing": 20}},
		{"Apprentice\n", map[string]int{}},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := otherRequirementSkillLevels(test.requirements)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For requirements %s, expected %v, but got %v", test.requirements, test.expected, result)
		}
	}
}

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

// TestExtractItemsFromIngredients is a unit test for extractItemsFromIngredients.
//func TestExtractItemsFromIngredients(t *testing.T) {
//	tests := []struct {
//		ingredientsString string
//		expectedItems     []Item
//	}{
//		{"HQ1: Antidote x6\nHQ2: Antidote x9\nHQ3: Antidote x12", []Item{{Name: "Antidote", Count: 6}, {Name: "Antidote", Count: 9}, {Name: "Antidote", Count: 12}}},
//		{"HQ1: Maple Shield +1\n", []Item{{Name: "Maple Shield +1", Count: 1}}},
//		{"HQ1: Angler's Hose\n", []Item{{Name: "Angler's Hose", Count: 1}}},
//		// Add more test cases as needed
//	}
//
//	for _, test := range tests {
//		result := extractItemsFromIngredients(test.ingredientsString)
//		if !reflect.DeepEqual(result, test.expectedItems) {
//			t.Errorf("For ingredients string %s, expected %v, but got %v", test.ingredientsString, test.expectedItems, result)
//		}
//	}
//}
