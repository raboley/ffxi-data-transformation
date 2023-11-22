package ffxi

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Recipe represents a crafting recipe.
type Recipe struct {
	CraftingType string `json:"crafting_type"`
	Crystal      string `json:"crystal"`
	// ... (other fields remain unchanged)
}

// CrystalData represents the data extracted for each crystal.
type CrystalData struct {
	Crystal string `json:"Crystal"`
}

// Item represents a crafting item.
type Item struct {
	Name  string `json:"Name"`
	Count int    `json:"Count"`
}

// ExtractCrystals processes the input JSON string and returns the resulting JSON string containing all crystals.
func ExtractCrystals(inputJSON string) (string, error) {
	var recipes []Recipe
	err := json.Unmarshal([]byte(inputJSON), &recipes)
	if err != nil {
		return "", err
	}

	// Create a slice of CrystalData objects
	var crystalData []CrystalData
	for _, recipe := range recipes {
		crystalData = append(crystalData, CrystalData{Crystal: recipe.Crystal})
	}

	// Marshal the CrystalData slice into JSON
	outputData, err := json.MarshalIndent(crystalData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(outputData), nil
}

// TransformRecipes is a placeholder function that currently calls ExtractCrystals.
// You can extend this function with additional transformations as needed.
func TransformRecipes(inputJSON string) (string, error) {
	// Placeholder implementation, currently calls ExtractCrystals
	return ExtractCrystals(inputJSON)
}

// extractCraftType extracts the craft type from the Text field using regex.
func extractCraftType(text string) string {
	re := regexp.MustCompile(`Guild Recipes: (\w+)`)
	match := re.FindStringSubmatch(text)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

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

// otherRequirementSkillLevels extracts skill levels from other requirements text.
func otherRequirementSkillLevels(requirements string) map[string]int {
	skillLevels := make(map[string]int)
	re := regexp.MustCompile(`(\w+)\((\d+)\)`)
	matches := re.FindAllStringSubmatch(requirements, -1)
	for _, match := range matches {
		craftType := match[1]
		level, _ := strconv.Atoi(match[2])
		skillLevels[craftType] = level
	}
	return skillLevels
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

// extractItemsFromIngredients extracts items and their counts from the ingredients string.
func extractItemsFromIngredients(ingredientsString string) []Item {
	var items []Item
	re := regexp.MustCompile(`HQ\d+: ([^\n]+?)(?: x(\d+))?`)
	matches := re.FindAllStringSubmatch(ingredientsString, -1)
	for _, match := range matches {
		name := strings.TrimSpace(match[1])
		var count int
		if len(match) == 3 && match[2] != "" {
			count, _ = strconv.Atoi(match[2])
		} else {
			count = 1
		}
		items = append(items, Item{Name: name, Count: count})
	}
	return items
}

// TestExtractItemsFromIngredients is a unit test for extractItemsFromIngredients.
func TestExtractItemsFromIngredients(t *testing.T) {
	tests := []struct {
		ingredientsString string
		expectedItems     []Item
	}{
		{"HQ1: Antidote x6\nHQ2: Antidote x9\nHQ3: Antidote x12", []Item{{Name: "Antidote", Count: 6}, {Name: "Antidote", Count: 9}, {Name: "Antidote", Count: 12}}},
		{"HQ1: Maple Shield +1\n", []Item{{Name: "Maple Shield +1", Count: 1}}},
		{"HQ1: Angler's Hose\n", []Item{{Name: "Angler's Hose", Count: 1}}},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := extractItemsFromIngredients(test.ingredientsString)
		if !reflect.DeepEqual(result, test.expectedItems) {
			t.Errorf("For ingredients string %s, expected %v, but got %v", test.ingredientsString, test.expectedItems, result)
		}
	}
}
