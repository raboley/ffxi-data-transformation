package recipe

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// TransformRecipes is a placeholder function that currently calls ExtractCrystals.
// You can extend this function with additional transformations as needed.
func TransformRecipes(inputJSON string) (string, error) {
	// Placeholder implementation, currently calls ExtractCrystals
	return ExtractCrystals(inputJSON)
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

// extractCraftType extracts the craft type from the Text field using regex.
func extractCraftType(text string) string {
	re := regexp.MustCompile(`Guild Recipes: (\w+)`)
	match := re.FindStringSubmatch(text)
	if len(match) == 2 {
		return match[1]
	}
	return ""
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
