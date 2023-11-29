package recipe

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// TransformRecipes processes the input JSON and returns the resulting JSON string.
func TransformRecipes(inputJSON string) ([]CraftingRecipe, error) {
	// Unmarshal the entire JSON data
	var recipes []CraftData
	err := json.Unmarshal([]byte(inputJSON), &recipes)
	if err != nil {
		return nil, err
	}

	// Create a slice of CraftingRecipe objects
	var craftingRecipes []CraftingRecipe
	for _, recipe := range recipes {
		// Extract craft type from the "Text" field using regex

		craftType := extractCraftType(recipe.Text)
		skillLevels := extractSkillLevels(recipe.LevelCap, craftType)
		otherSkillLevels := otherRequirementSkillLevels(recipe.OtherRequirements)
		combinedSkillLevels := combineSkillLevels(skillLevels, otherSkillLevels)

		sortedSkills := sortSkillsHighestFirst(combinedSkillLevels)
		realMainCraftType := sortedSkills[0]

		items := extractRequiredItems(recipe.SynthOrDesynth)
		name := determineCraftName(combinedSkillLevels, items, recipe.RecipeName)
		recipeQuantity := extractRecipeQuantity(recipe.RecipeName)
		standardResult := ResultsIncludingHighQuality{
			Name:             recipe.RecipeItem,
			Count:            recipeQuantity,
			HighQualityLevel: 0,
		}
		highQualityResults, _ := extractHighQualityResults(recipe.Ingredients)
		allResults := append(highQualityResults, standardResult)
		requiredTools := extractToolRequirement(recipe.OtherRequirements)

		// Create CraftingRecipe object
		craftingRecipes = append(craftingRecipes, CraftingRecipe{
			Result:             recipe.RecipeItem,
			Crystal:            recipe.Crystal,
			MainCraft:          realMainCraftType,
			SkillLevels:        combinedSkillLevels,
			RequiredItems:      items,
			Name:               name,
			AllPossibleResults: allResults,
			RequiredTools:      requiredTools,
		})
	}

	return craftingRecipes, nil

}

func extractRecipeQuantity(itemName string) int {
	re := regexp.MustCompile(` x(\d+)$`)
	match := re.FindStringSubmatch(itemName)

	if len(match) > 1 {
		quantity, _ := strconv.Atoi(match[1])
		return quantity
	}

	return 1
}

type ResultsIncludingHighQuality struct {
	Name             string `json:"Name"`
	Count            int    `json:"Count"`
	HighQualityLevel int    `json:"HighQualityLevel"`
}

func extractSkillLevels(levelCap string, craftType string) map[string]int {
	intLevelCap, _ := strconv.Atoi(levelCap)
	baseCraftLevel := map[string]int{
		craftType: intLevelCap,
	}
	return baseCraftLevel
}

//// extractCraftType extracts the craft type from the "Text" field using a simple regex.
//func extractCraftType(text string) string {
//	// Use a simple regex to extract craft type
//	re := regexp.MustCompile(`Guild Recipes: ([a-zA-Z]+)`)
//	match := re.FindStringSubmatch(text)
//	if len(match) > 1 {
//		return strings.TrimSpace(match[1])
//	}
//	return ""
//}

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

// combineSkillLevels combines two maps of skill levels.
func combineSkillLevels(dest, src map[string]int) map[string]int {
	combined := make(map[string]int)

	// Copy the destination map to the combined map
	for k, v := range dest {
		combined[k] = v
	}

	// Add or update entries from the source map
	for k, v := range src {
		if existing, ok := combined[k]; ok {
			// Choose how to handle conflicts (e.g., sum values)
			combined[k] = existing + v
		} else {
			combined[k] = v
		}
	}

	return combined
}

// extractRequiredItems extracts required items from a string.
func extractRequiredItems(itemsString string) []Item {
	// Split the string by commas
	items := strings.Split(itemsString, ",")

	var requiredItems []Item

	// Regular expression to match quantity indicators like "x3"
	re := regexp.MustCompile(`x(\d+)`)

	for _, item := range items {
		// Trim spaces from the item
		item = strings.TrimSpace(item)

		// Check if there's a quantity indicator
		matches := re.FindStringSubmatch(item)
		if len(matches) > 1 {
			count, err := strconv.Atoi(matches[1])
			if err != nil {
				// Handle error
				continue
			}
			// Extract the item name without the quantity indicator
			name := strings.TrimSpace(re.ReplaceAllString(item, ""))
			requiredItems = append(requiredItems, Item{Name: name, Count: count})
		} else {
			// No quantity indicator, assume count is 1
			requiredItems = append(requiredItems, Item{Name: item, Count: 1})
		}
	}

	return requiredItems
}

func extractHighQualityResults(ingredients string) ([]ResultsIncludingHighQuality, error) {
	var results []ResultsIncludingHighQuality

	itemLines := strings.Split(ingredients, "\n")

	for _, line := range itemLines {
		fmt.Printf("line: %s\n", line)
		// Use the updated regular expression to capture item details
		re := regexp.MustCompile(`HQ(\d+): (.*?)(?: x(\d+))?$`)

		match := re.FindStringSubmatch(line)
		fmt.Printf("match: %v\n", strings.Join(match, ", "))
		if len(match) > 0 {
			hqLevel, _ := strconv.Atoi(match[1])
			itemName := match[2]
			quantity := 1

			if match[3] != "" {
				quantity, _ = strconv.Atoi(match[3])
			}

			fmt.Printf("itemName: %s\n", itemName)
			fmt.Printf("itemQuantity: %s\n", match[3])

			results = append(results, ResultsIncludingHighQuality{
				Name:             itemName,
				Count:            quantity,
				HighQualityLevel: hqLevel,
			})
		}
	}

	return results, nil
}

// determineCraftName determines the name of the craft based on the result and required items.
func determineCraftName(skillLevels map[string]int, requiredItems []Item, result string) string {

	craftTypes := sortSkillsHighestFirst(skillLevels)

	// Iterate over craft types in descending order
	var skills []string
	for _, craftType := range craftTypes {
		level := skillLevels[craftType]
		skills = append(skills, fmt.Sprintf("%s-%d", craftType, level))
		fmt.Printf("%s-%d\n", craftType, level)
	}

	var items []string
	for _, item := range requiredItems {
		items = append(items, fmt.Sprintf("%d-%s", item.Count, item.Name))
	}

	// Combine the result and items to form the craft name
	craftName := fmt.Sprintf("%s-%s-From-%s", strings.Join(skills, "-"), result, strings.Join(items, ", "))
	return craftName
}

func sortSkillsHighestFirst(skillLevels map[string]int) []string {
	// Create a slice of craft types
	var craftTypes []string
	for craftType := range skillLevels {
		craftTypes = append(craftTypes, craftType)
	}

	// Sort the craft types based on values in descending order
	sort.Slice(craftTypes, func(i, j int) bool {
		return skillLevels[craftTypes[i]] > skillLevels[craftTypes[j]]
	})
	return craftTypes
}

// extractToolRequirement extracts the tool requirement from the given text.
// It returns the tool name or an empty string if no tool requirement is found.
func extractToolRequirement(text string) string {
	re := regexp.MustCompile(`Tool: (.*)`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// CraftingRecipe represents the data extracted for each craft.
type CraftingRecipe struct {
	Crystal            string                        `json:"Crystal"`
	RequiredItems      []Item                        `json:"RequiredItems"`
	SkillLevels        map[string]int                `json:"SkillLevels"`
	Result             string                        `json:"Result"`
	Name               string                        `json:"Name"`
	MainCraft          string                        `json:"MainCraft"`
	AllPossibleResults []ResultsIncludingHighQuality `json:"AllPossibleResults"`
	RequiredTools      string                        `json:"RequiredTools"`
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

// CraftData represents a crafting recipe.
type CraftData struct {
	Text                       string `json:"Text"`
	RecipeName                 string `json:"recipe_name"`
	GuildRecipesWoodworkingURL string `json:"Guild_Recipes_Woodworking_URL"`
	RecipeItem                 string `json:"recipe_item"`
	LevelCap                   string `json:"level_cap"`
	OtherRequirements          string `json:"other_requirements"`
	Crystal                    string `json:"crystal"`
	SynthOrDesynth             string `json:"synth_or_desynth"`
	Ingredients                string `json:"ingredients"`
	Something                  string `json:"something"`
	Ingredient1                string `json:"ingredient_1"`
	Ingredient2Link            string `json:"ingredient_2_link"`
	Ingredient2                string `json:"ingredient_2"`
	Ingredient3Link            string `json:"ingredient_3_link"`
	Ingredient3                string `json:"ingredient_3"`
	Ingredient4Link            string `json:"ingredient_4_link"`
	Ingredient4                string `json:"ingredient_4"`
	HQResults                  string `json:"hq_results"`
	Field15                    string `json:"Field15"`
	HQ1                        string `json:"hq1"`
	HQ2                        string `json:"hq2"`
	HQ3                        string `json:"hq3"`
}
