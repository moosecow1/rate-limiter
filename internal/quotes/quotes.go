package quotes

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

const quotesFilePath = "quotes.json"

var quotes []Quote

func init() {
	var err error
	quotes, err = LoadQuotes()

	if err != nil {
		fmt.Errorf("Failed to load quotes: %w", err)
		return
	}
}

func LoadQuotes() ([]Quote, error) {

	content, err := os.ReadFile(quotesFilePath)

	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %w", err)
	}

	var quotes []Quote

	err = json.Unmarshal(content, &quotes)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %w", err)
	}

	return quotes, nil

}

func PickQuote() (Quote, error) {

	/*
		quotes, err := LoadQuotes()

		if err != nil {
			return Quote{}, fmt.Errorf("Error loading quotes: %w", err)
		}
	*/

	if len(quotes) < 1 {
		return Quote{}, fmt.Errorf("There are no quotes.")
	}

	num := rand.Intn(len(quotes))

	return quotes[num], nil

}

func PickQuoteFromAuthor(author string) (Quote, error) {

	/*
		quotes, err := LoadQuotes()

		if err != nil {
			return Quote{}, fmt.Errorf("Error loading quotes: %w", err)
		}
	*/

	filtered := quotes[:0]

	log.Printf("%t %s", author == "", author)

	if author == "" {
		filtered = quotes
	} else {
		for i := len(quotes) - 1; i >= 0; i-- {
			if quotes[i].Author == author {
				filtered = append(filtered, quotes[i])
			}
		}
	}

	if len(filtered) < 1 {
		return Quote{}, fmt.Errorf("There are no quotes.")
	}

	num := rand.Intn(len(filtered))

	return filtered[num], nil

}
