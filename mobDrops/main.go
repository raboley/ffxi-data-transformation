package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type ItemInfo struct {
	ItemName string `json:"Name"`
	NPC      string `json:"NPC"`
	Zone     string `json:"Zone"`
	Count    string `json:"Count"`
	Chance   string `json:"Chance"`
	PageURL  string `json:"Page_URL"`
}

type ItemDrop struct {
	Name           string  `json:"Name"`
	Percent        float64 `json:"Percent"`
	AmountDropped  int     `json:"AmountDropped"`
	AmountDefeated int     `json:"AmountDefeated"`
}

type MobInfo struct {
	Name       string `json:"Name"`
	LevelRange struct {
		Min int `json:"Min"`
		Max int `json:"Max"`
	} `json:"LevelRange"`
	ZoneName  string     `json:"ZoneName"`
	ItemDrops []ItemDrop `json:"ItemDrops"`
}

func main() {
	// Load item drop info
	itemInfo, err := loadItemInfo("all_mobs_nineth.json")
	if err != nil {
		log.Fatal(err)
	}

	// Load and update mob info for each zone
	zones := []string{
		"Valkurm_Dunes",
		//"Jugner_Forest",
	} // Add more zones as needed

	for _, zone := range zones {
		mobInfo, err := loadMobInfo(zone + ".json")
		if err != nil {
			log.Fatal(err)
		}

		updateDropChances(mobInfo, itemInfo)

		err = writeMobInfo(zone+"_output.json", mobInfo)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func loadItemInfo(filename string) ([]ItemInfo, error) {
	// Read the content of the file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Use regex to replace single quotes with double quotes not preceded by backslash
	re := regexp.MustCompile(`(?<!\\)'`)
	doubleQuotedJSON := re.ReplaceAllString(string(fileContent), `"`)

	// Unmarshal the JSON
	var itemInfo []ItemInfo
	err = json.Unmarshal([]byte(doubleQuotedJSON), &itemInfo)
	if err != nil {
		log.Printf("Error unmarshalling JSON from file %s: %v", filename, err)
		return nil, err
	}

	return itemInfo, nil
}

func loadMobInfo(filename string) ([]MobInfo, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var mobInfo []MobInfo
	err = json.Unmarshal(fileContent, &mobInfo)
	if err != nil {
		return nil, err
	}

	return mobInfo, nil
}

func updateDropChances(mobInfo []MobInfo, itemInfo []ItemInfo) {
	for i, mob := range mobInfo {
		var updatedItemDrops []ItemDrop

		for _, item := range mob.ItemDrops {
			// Find the corresponding item info
			var info ItemInfo
			for _, itemInfoEntry := range itemInfo {
				if strings.EqualFold(item.Name, itemInfoEntry.ItemName) && strings.EqualFold(mob.Name, itemInfoEntry.NPC) && strings.EqualFold(mob.ZoneName, itemInfoEntry.Zone) {
					info = itemInfoEntry
					break
				}
			}

			// If item info is found, update drop chances
			if info.ItemName != "" {
				percent, err := parsePercent(info.Chance)
				if err != nil {
					log.Printf("Error parsing percent for item %s: %v", item.Name, err)
					percent = 100.0 // Default to 100% if parsing fails
				}

				updatedItemDrops = append(updatedItemDrops, ItemDrop{
					Name:           item.Name,
					Percent:        percent,
					AmountDropped:  item.AmountDropped,
					AmountDefeated: item.AmountDefeated,
				})
			} else {
				// If item info is not found, set default drop chances
				updatedItemDrops = append(updatedItemDrops, ItemDrop{
					Name:           item.Name,
					Percent:        100.0,
					AmountDropped:  0,
					AmountDefeated: 0,
				})
			}
		}

		// Update the mob's ItemDrops field
		mobInfo[i].ItemDrops = updatedItemDrops
	}
}

func parsePercent(percentStr string) (float64, error) {
	var percent float64
	_, err := fmt.Sscanf(percentStr, "%f%%", &percent)
	if err != nil {
		return 0, err
	}

	return percent, nil
}

func writeMobInfo(filename string, mobInfo []MobInfo) error {
	// Marshal the updated data
	updatedData, err := json.MarshalIndent(mobInfo, "", "  ")
	if err != nil {
		return err
	}

	// Write the updated data back to the file
	err = ioutil.WriteFile(filename, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
