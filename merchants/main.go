package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Merchant struct {
	Merchant   string `json:"merchant"`
	Type       string `json:"type"`
	GoodsPrice string `json:"goodsPrice"`
	Location   string `json:"location"`
}

type ItemInfo struct {
	Name            string
	MinPrice        int
	MaxPrice        int
	RankRequirement string
}

type MerchantInfo struct {
	Name  string
	Items []ItemInfo
	Zone  string
}

func extractMerchantInfo(jsonData []byte) ([]MerchantInfo, error) {
	var merchants []Merchant
	if err := json.Unmarshal(jsonData, &merchants); err != nil {
		return nil, err
	}

	var merchantInfoList []MerchantInfo
	for _, merchant := range merchants {
		goodsList, err := extractGoodsAndPrices(merchant.GoodsPrice)
		if err != nil {
			return nil, err
		}
		zone, err := extractZone(merchant.Location)
		if err != nil {
			return nil, err
		}

		merchantInfoList = append(merchantInfoList, MerchantInfo{
			Name:  merchant.Merchant,
			Items: goodsList,
			Zone:  zone,
		})
	}

	return merchantInfoList, nil
}

func extractGoodsAndPrices(goodsPrice string) ([]ItemInfo, error) {
	re := regexp.MustCompile(`([^\d]+) (\d+)(?:-(\d+))? gil`)
	matches := re.FindAllStringSubmatch(goodsPrice, -1)

	var items []ItemInfo
	for _, match := range matches {
		itemName := strings.TrimSpace(match[1])
		// Handle cases where the item name contains ranking information
		if strings.Contains(itemName, "\n") {
			// Extract rank information from the item text
			rankMatch := regexp.MustCompile(`(\w+ place)`)
			rankSubmatch := rankMatch.FindStringSubmatch(itemName)
			if len(rankSubmatch) > 1 {
				rank := rankSubmatch[1]
				// Remove rank information from the item name
				itemName = strings.TrimSpace(strings.Replace(itemName, rank, "", 1))
				// Create an ItemInfo struct and include rank information
				minPrice, err := strconv.Atoi(match[2])
				if err != nil {
					return nil, err
				}
				var maxPrice int
				if match[3] != "" {
					maxPrice, err = strconv.Atoi(match[3])
					if err != nil {
						return nil, err
					}
				} else {
					maxPrice = minPrice
				}
				item := ItemInfo{
					Name:            itemName,
					MinPrice:        minPrice,
					MaxPrice:        maxPrice,
					RankRequirement: rank,
				}
				items = append(items, item)
			}
		} else {
			minPrice, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, err
			}
			// Check if there's a specified max price
			var maxPrice int
			if match[3] != "" {
				maxPrice, err = strconv.Atoi(match[3])
				if err != nil {
					return nil, err
				}
			} else {
				// If no max price is specified, set it equal to the min price
				maxPrice = minPrice
			}
			// Create an ItemInfo struct without rank information
			item := ItemInfo{
				Name:     itemName,
				MinPrice: minPrice,
				MaxPrice: maxPrice,
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func extractZone(location string) (string, error) {
	re := regexp.MustCompile(`([^\(]+) \([^\)]+\)`)
	match := re.FindStringSubmatch(location)
	if len(match) > 1 {
		zone := strings.TrimSpace(match[1])
		return sanitizeZoneName(zone), nil
	}
	return "", fmt.Errorf("unable to extract zone from location: %s", location)
}

func sanitizeZoneName(zone string) string {
	// Split the zone name into words
	words := strings.Fields(zone)

	// Check for cardinal directions or the word "port" and move them to the front
	for i, word := range words {
		lowercaseWord := strings.ToLower(word)
		if lowercaseWord == "north" || lowercaseWord == "south" || lowercaseWord == "east" || lowercaseWord == "west" || lowercaseWord == "port" {
			// Move the word to the front
			words = append([]string{word}, append(words[:i], words[i+1:]...)...)
			break
		}
	}

	// Join the words back together
	sanitized := strings.Join(words, "_")
	// Remove special characters and apostrophes
	sanitized = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(sanitized, "")

	return sanitized
}

func writeMerchantFiles(merchantInfoList []MerchantInfo) error {
	// Create a directory for merchants
	err := os.MkdirAll("merchants", os.ModePerm)
	if err != nil {
		return err
	}

	// Group merchantInfo by zone
	zoneMerchants := make(map[string][]MerchantInfo)
	for _, info := range merchantInfoList {
		zoneMerchants[info.Zone] = append(zoneMerchants[info.Zone], info)
	}

	// Write files for each zone
	for zone, merchants := range zoneMerchants {
		// Convert merchants to JSON
		jsonData, err := json.MarshalIndent(merchants, "", "  ")
		if err != nil {
			return err
		}

		// Write the file for each zone
		filePath := filepath.Join("merchants", fmt.Sprintf("%s.json", zone))
		err = ioutil.WriteFile(filePath, jsonData, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Read JSON data from the input file
	jsonData, err := ioutil.ReadFile("input.json")
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	merchantInfoList, err := extractMerchantInfo(jsonData)
	if err != nil {
		fmt.Println("Error extracting merchant info:", err)
		return
	}

	err = writeMerchantFiles(merchantInfoList)
	if err != nil {
		fmt.Println("Error writing merchant files:", err)
		return
	}

	fmt.Println("Name files written successfully.")
}
