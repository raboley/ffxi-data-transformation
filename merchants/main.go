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
	ItemName string
	MinPrice int
	MaxPrice int
}

type MerchantInfo struct {
	Merchant string
	Items    []ItemInfo
	Zone     string
	Rank     string // New field to store the ranking information
}

func extractMerchantInfo(jsonData []byte) ([]MerchantInfo, error) {
	var merchants []Merchant
	if err := json.Unmarshal(jsonData, &merchants); err != nil {
		return nil, err
	}

	var merchantInfoList []MerchantInfo
	for _, merchant := range merchants {
		goodsList, prices := extractGoodsAndPrices(merchant.GoodsPrice)
		items := make([]ItemInfo, len(goodsList))
		for i, item := range goodsList {
			items[i] = ItemInfo{
				ItemName: strings.TrimSpace(item),
				MinPrice: prices[i*2],
				MaxPrice: prices[i*2+1],
			}
		}

		zone, rank := extractZoneAndRank(merchant.Location)

		merchantInfoList = append(merchantInfoList, MerchantInfo{
			Merchant: merchant.Merchant,
			Items:    items,
			Zone:     zone,
			Rank:     rank,
		})
	}

	return merchantInfoList, nil
}

func extractGoodsAndPrices(goodsPrice string) ([]string, []int) {
	re := regexp.MustCompile(`([^\d]+) (\d+)-(\d+) gil`)
	matches := re.FindAllStringSubmatch(goodsPrice, -1)

	var goodsList []string
	var prices []int
	for _, match := range matches {
		itemName := strings.TrimSpace(match[1])
		// Handle cases where the item name contains ranking information
		if strings.Contains(itemName, "\n") {
			itemName = strings.TrimSpace(strings.Split(itemName, "\n")[1])
		}
		goodsList = append(goodsList, itemName)
		minPrice, _ := strconv.Atoi(match[2])
		maxPrice, _ := strconv.Atoi(match[3])
		prices = append(prices, minPrice, maxPrice)
	}

	return goodsList, prices
}

func extractZoneAndRank(location string) (string, string) {
	re := regexp.MustCompile(`([^\(]+) \(([^\)]+)\)`)
	match := re.FindStringSubmatch(location)
	if len(match) > 2 {
		return strings.TrimSpace(match[1]), strings.TrimSpace(match[2])
	}
	return "", ""
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

	fmt.Println("Merchant files written successfully.")
}
