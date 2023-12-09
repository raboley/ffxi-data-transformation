package main

import (
	"encoding/json"
	"ffxi/recipe"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//fileName := "Goldsmithing Guild Recipes _Synthesis_ - Final Fantasy XI - somepage.com.json"
	//fileName := "clothcraft Guild Recipes _Synthesis_ - Final Fantasy XI - somepage.com.json"
	//fileName := "Leatherworking Guild Recipes _Synthesis_ - Final Fantasy XI - somepage.com.json"
	fileName := "smithing and woodworking Guild Recipes _Synthesis_ - Final Fantasy XI - somepage.com"
	filePath := fmt.Sprintf("/Users/russellboley/Documents/%s", fileName)
	inputJSON, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	craftingRecipes, err := recipe.TransformRecipes(string(inputJSON))
	if err != nil {
		log.Fatal(err)
	}

	allRecipeFilePath := "/Users/russellboley/git/FantasyAi/FinalFantasyData/CraftingRecipes/not_ready_for_yet"
	os.Mkdir(allRecipeFilePath, 0777)
	for _, recipe := range craftingRecipes {
		recipeJSON, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fileName, err := getShortFileName(recipe.Name)
		if err != nil {
			log.Fatal(err)
		}
		filePath := fmt.Sprintf("%s/%s.json", allRecipeFilePath, fileName)

		_, err = os.Stat(filePath)
		if err == nil {
			continue

			//filePath = addCharacterBeforeExtension(filePath, '2')
			//log.Fatalf("File at path %s already exists", filePath)
		}

		err = os.WriteFile(filePath, []byte(recipeJSON), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Marshal the CraftingRecipe slice into JSON
	outputJSON, err := json.MarshalIndent(craftingRecipes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("all_craft.json", []byte(outputJSON), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Extract item names from the transformed JSON
	var items []string
	var recipes []recipe.CraftingRecipe
	err = json.Unmarshal([]byte(outputJSON), &recipes)
	if err != nil {
		log.Fatal(err)
	}

	uniqueItems := make(map[string]struct{})

	for _, r := range recipes {
		for _, i := range r.RequiredItems {
			if _, ok := uniqueItems[i.Name]; !ok {
				items = append(items, i.Name)
				// Add the item name to the map
				uniqueItems[i.Name] = struct{}{}
			}
		}
	}

	// Write item names to a file
	err = os.WriteFile("item_names.txt", []byte(strings.Join(items, "\n")), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transformation complete. Output written to output.json")
}

func getShortFileName(fileName string) (string, error) {
	// Extract everything before "-From" (non-inclusive)
	index := strings.Index(fileName, ",")
	if index == -1 {
		// file name is probably pretty short
		return fileName, nil
	}

	resultFilename := strings.TrimSpace(fileName[:index])

	// Check if the resulting filename is empty
	if resultFilename == "" {
		log.Fatal("Invalid resulting filename")
	}

	// Display the resulting filename
	fmt.Println("Resulting filename:", resultFilename)

	return resultFilename, nil
}

// addCharacterBeforeExtension adds a character before the file extension.
func addCharacterBeforeExtension(filePath string, char rune) string {
	ext := filepath.Ext(filePath)
	if ext == "" {
		// No extension found
		return filePath
	}

	// Find the position of the last dot before the extension
	dotIndex := strings.LastIndex(filePath, ext)

	// Insert the character before the dot
	newFilePath := filePath[:dotIndex] + string(char) + filePath[dotIndex:]

	return newFilePath
}
