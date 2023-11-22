package ffxi

import (
	"regexp"
	"strconv"
	"strings"
)

// Item represents an item and its count.
type Item struct {
	Name  string
	Count int
}

// extractItemsFromIngredients extracts items and their counts from the ingredients string.
func extractItemsFromIngredients(ingredientsString string) []Item {
	var items []Item

	// Regular expression to match item names and counts
	re := regexp.MustCompile(`HQ\d+: ([^\n]+?)(?: x(\d+))?`)

	// Find all matches in the ingredients string
	matches := re.FindAllStringSubmatch(ingredientsString, -1)

	for _, match := range matches {
		name := strings.TrimSpace(match[1])

		// Check if the count is present, otherwise assume count is 1
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
