package distributors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Mouser represents a distributor that implements the Distributor interface
type Mouser struct {
	BaseURL string
	APIKey  string
}

// Initialize sets the API key for Mouser
func (m *Mouser) Initialize(params ...string) error {
	if len(params) < 1 {
		return fmt.Errorf("Mouser initialization requires an API key")
	}
	apiKey := params[0]

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	m.APIKey = apiKey
	m.BaseURL = "https://api.mouser.com/api/v1/"
	return nil
}

// QueryPartNumber queries Mouser's API for a given part number
func (m *Mouser) QueryPartNumber(partNumber string) (PartInfo, error) {
	// Handle a empty PN
	if partNumber == "" {
		return PartInfo{}, fmt.Errorf("part number cannot be empty")
	}

	url := fmt.Sprintf("%s/search/partnumber?apiKey=%s", m.BaseURL, m.APIKey)

	// Create the JSON request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"SearchByPartRequest": map[string]interface{}{
			"mouserPartNumber": partNumber,
		},
	})
	if err != nil {
		return PartInfo{}, fmt.Errorf("failed to marshal JSON request body: %v", err)
	}

	// Make the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return PartInfo{}, fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return PartInfo{}, fmt.Errorf("Mouser API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		Errors []struct {
			Id                    int    `json: Id`
			Code                  string `json:"Code"`
			Message               string `json:Message`
			ResourceKey           string `json:ResourceKey`
			ResourceFormatString  string `json: ResourceFormatString`
			ResourceFormatString2 string `json:ResourceFormatString2`
			PropertyName          string `json:PropertyName`
		} `json:"Errors`
		SearchResults struct {
			NumberOfResult int `json:"NumberOfResult"`
			Parts          []struct {
				PartNumber             string `json:"MouserPartNumber"`
				ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
				Manufacturer           string `json:"Manufacturer"`
				Description            string `json:"Description"`
				PriceBreaks            []struct {
					Quantity int    `json:"Quantity`
					Price    string `json:"Price`
					Currency string `json:"Currency`
				}
				AvailabilityInStock string  `json:"AvailabilityInStock"`
				UnitPrice           float64 `json:"unitprice"`
				ProductURL          string  `json:"ProductDetailUrl"`
				DataSheetURL        string  `json:"DataSheetUrl"`
			} `json:"Parts"`
		} `json:"SearchResults`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return PartInfo{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if len(result.Errors) != 0 {
		return PartInfo{}, fmt.Errorf("encountered a Mouser error regarding %s", result.Errors[0].PropertyName)
	}

	if result.SearchResults.NumberOfResult == 0 {
		return PartInfo{}, fmt.Errorf("no parts found for part number: %s", partNumber)
	}

	priceStr := strings.TrimPrefix(result.SearchResults.Parts[0].PriceBreaks[0].Price, "$")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return PartInfo{}, fmt.Errorf("Error converting unit price to float64")
	}

	availability, err := strconv.Atoi(result.SearchResults.Parts[0].AvailabilityInStock)
	if err != nil {
		return PartInfo{}, fmt.Errorf("Error converting availability to int")
	}
	// Convert to the unified format
	partInfo := PartInfo{
		PartNumber:             result.SearchResults.Parts[0].PartNumber,
		ManufacturerPartNumber: result.SearchResults.Parts[0].ManufacturerPartNumber,
		Manufacturer:           result.SearchResults.Parts[0].Manufacturer,
		Description:            result.SearchResults.Parts[0].Description,
		UnitPrice:              price,
		Availability:           availability,
		ProductURL:             result.SearchResults.Parts[0].ProductURL,
		DataSheetURL:           result.SearchResults.Parts[0].DataSheetURL,
	}

	return partInfo, nil
}
