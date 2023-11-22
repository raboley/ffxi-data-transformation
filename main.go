package main

import (
	"ffxi/recipe"
	"fmt"
	"log"
	"os"
)

func main() {
	inputJSON, err := os.ReadFile("/Users/russellboley/Documents/all_crafting.json")
	if err != nil {
		log.Fatal(err)
	}

	outputJSON, err := recipe.TransformRecipes(string(inputJSON))
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("all_craft.json", []byte(outputJSON), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transformation complete. Output written to output.json")
}
