package recipe

import (
	"fmt"
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

func TestExtractSkillLevels(t *testing.T) {
	tests := []struct {
		levelCap    string
		craftType   string
		expectedMap map[string]int
	}{
		{"10", "Woodworking", map[string]int{"Woodworking": 10}},
		{"15", "Smithing", map[string]int{"Smithing": 15}},
		{"5", "Alchemy", map[string]int{"Alchemy": 5}},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := extractSkillLevels(test.levelCap, test.craftType)
		if !reflect.DeepEqual(result, test.expectedMap) {
			t.Errorf("For CraftType %s and LevelCap %s, expected %v, but got %v", test.craftType, test.levelCap, test.expectedMap, result)
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

func TestExtractHighQualityResults(t *testing.T) {
	tests := []struct {
		ingredientsString string
		expectedItems     []ResultsIncludingHighQuality
	}{
		{
			"HQ1: Antidote x6\nHQ2: Antidote x9\nHQ3: Antidote x12",
			[]ResultsIncludingHighQuality{
				{Name: "Antidote", Count: 6, HighQualityLevel: 1},
				{Name: "Antidote", Count: 9, HighQualityLevel: 2},
				{Name: "Antidote", Count: 12, HighQualityLevel: 3},
			},
		},
		{
			"HQ1: Maple Shield +1\n",
			[]ResultsIncludingHighQuality{
				{Name: "Maple Shield +1", Count: 1, HighQualityLevel: 1},
			},
		},
		{
			"HQ1: Angler's Hose\n",
			[]ResultsIncludingHighQuality{
				{Name: "Angler's Hose", Count: 1, HighQualityLevel: 1},
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result, _ := extractHighQualityResults(test.ingredientsString)
		if !reflect.DeepEqual(result, test.expectedItems) {
			t.Errorf("For ingredients string %s, expected %v, but got %v", test.ingredientsString, test.expectedItems, result)
		}
	}
}

func TestExtractRequiredItems(t *testing.T) {
	tests := []struct {
		inputString  string
		expectedList []Item
	}{
		{"Wijnruit x3, San d'Orian Grape x3, Distilled Water, Triturator",
			[]Item{
				{Name: "Wijnruit", Count: 3},
				{Name: "San d'Orian Grape", Count: 3},
				{Name: "Distilled Water", Count: 1},
				{Name: "Triturator", Count: 1},
			},
		},
		{"Iron Ingot x2, Fire Crystal", []Item{
			{Name: "Iron Ingot", Count: 2},
			{Name: "Fire Crystal", Count: 1},
		}},
	}

	for _, test := range tests {
		result := extractRequiredItems(test.inputString)
		if !reflect.DeepEqual(result, test.expectedList) {
			t.Errorf("For input %s, expected %v, but got %v", test.inputString, test.expectedList, result)
		}
	}
}

func TestDetermineCraftName(t *testing.T) {
	tests := []struct {
		recipe       CraftingRecipe
		expectedName string
	}{
		{
			CraftingRecipe{
				Result: "Ash Lumber",
				RequiredItems: []Item{
					{Name: "Ash Log", Count: 1},
				},
				SkillLevels: map[string]int{
					"Woodworking": 7,
				},
			},
			"Woodworking-7-Ash Lumber-From-1-Ash Log",
		},
		{
			CraftingRecipe{
				Result: "Copper Ingot",
				RequiredItems: []Item{
					{Name: "Copper Ore", Count: 2},
					{Name: "Fire Crystal", Count: 1},
				},
				SkillLevels: map[string]int{
					"Goldsmithing": 2,
					"Smithing":     14,
				},
				MainCraft: "Goldsmithing",
			},
			"Smithing-14-Goldsmithing-2-Copper Ingot-From-2-Copper Ore, 1-Fire Crystal",
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := determineCraftName(test.recipe.SkillLevels, test.recipe.RequiredItems, test.recipe.Result)
		if result != test.expectedName {
			t.Errorf("For recipe %+v, expected %s, but got %s", test.recipe, test.expectedName, result)
		}
	}
}

func TestExtractRecipeQuantity(t *testing.T) {
	testCases := []struct {
		itemName string
		expected int
	}{
		{"Antidote x3", 3},
		{"Maple Shield +1", 1},
		{"Elixir", 1}, // No quantity specified
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Item: %s", testCase.itemName), func(t *testing.T) {
			result := extractRecipeQuantity(testCase.itemName)

			if result != testCase.expected {
				t.Errorf("Expected quantity %d, but got %d", testCase.expected, result)
			}
		})
	}
}

func TestExtractToolRequirement(t *testing.T) {
	testCases := []struct {
		text     string
		expected string
	}{
		{"Tool: Leather Ensorcellment", "Leather Ensorcellment"},
		{"Tool: Alchemic Ensorcellment", "Alchemic Ensorcellment"},
		{"Tool: Smithing Implements", "Smithing Implements"},
		{"Other Requirement", ""}, // No tool requirement
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Text: %s", testCase.text), func(t *testing.T) {
			result := extractToolRequirement(testCase.text)

			if result != testCase.expected {
				t.Errorf("Expected tool requirement %s, but got %s", testCase.expected, result)
			}
		})
	}
}

// TestExtractItemsFromIngredients is a unit test for extractHighQualityResults.
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
//		result := extractHighQualityResults(test.ingredientsString)
//		if !reflect.DeepEqual(result, test.expectedItems) {
//			t.Errorf("For ingredients string %s, expected %v, but got %v", test.ingredientsString, test.expectedItems, result)
//		}
//	}
//}
