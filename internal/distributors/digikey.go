package distributors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DigiKey struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	OAuthToken   string
	Locale       string
	Currency     string
}

// Initialize sets the client ID and client secret for DigiKey
func (d *DigiKey) Initialize(params ...string) error {
	if len(params) < 2 {
		return fmt.Errorf("DigiKey initialization requires client ID and client secret")
	}

	clientID := params[0]
	clientSecret := params[1]

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("client ID and client secret cannot be empty")
	}

	d.ClientID = clientID
	d.ClientSecret = clientSecret
	d.Locale = "US"
	d.Currency = "USD"
	d.BaseURL = "https://api.digikey.com"

	// Authenticate and obtain an access token
	accessToken, err := d.authenticate()
	if err != nil {
		return fmt.Errorf("failed to authenticate with DigiKey: %v", err)
	}

	d.OAuthToken = accessToken
	return nil
}

// QueryPartNumber queries DigiKey's API for a given part number
func (d *DigiKey) QueryPartNumber(partNumber string) (PartInfo, error) {
	if partNumber == "" {
		return PartInfo{}, fmt.Errorf("part number cannot be empty")
	}

	api_url := fmt.Sprintf("%s/products/v4/search/%s/productdetails", d.BaseURL, url.QueryEscape(partNumber))

	// Create a new GET request
	req, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		return PartInfo{}, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("X-DIGIKEY-Locale-Site", d.Locale)
	req.Header.Set("X-DIGIKEY-Locale-Currency", d.Currency)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.OAuthToken))
	req.Header.Set("X-DIGIKEY-Client-Id", d.ClientID)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PartInfo{}, fmt.Errorf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return PartInfo{}, fmt.Errorf("Digikey API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		Product struct {
			ProductVariations []struct {
				ProductNumber string `json:"DigiKeyProductNumber"`
			} `json:"ProductVariations"`
			PartNumber             string `json:"MouserPartNumber"`
			ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
			Manufacturer           struct {
				Name string `json:"Name"`
			} `json:"Manufacturer"`
			Description struct {
				ProductDescription  string `json:"ProductDescription"`
				DetailedDescription string `json:"DetailedDescription"`
			} `json:"Description"`
			PriceBreaks []struct {
				Quantity int    `json:"Quantity`
				Price    string `json:"Price`
				Currency string `json:"Currency`
			}
			QuanitityAvailable int     `json:"QuanitityAvailable"`
			UnitPrice          float64 `json:UnitPrice"`
			ProductURL         string  `json:"ProductDetailUrl"`
			DataSheetURL       string  `json:"DataSheetUrl"`
		} `json:"Product`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return PartInfo{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}
	// Convert to the unified format
	partInfo := PartInfo{
		PartNumber:             result.Product.ProductVariations[0].ProductNumber,
		ManufacturerPartNumber: result.Product.ManufacturerPartNumber,
		Manufacturer:           result.Product.Manufacturer.Name,
		Description:            result.Product.Description.ProductDescription,
		UnitPrice:              result.Product.UnitPrice,
		Availability:           result.Product.QuanitityAvailable,
		ProductURL:             result.Product.ProductURL,
		DataSheetURL:           result.Product.DataSheetURL,
	}

	return partInfo, nil
}

// authenticate handles the logic for obtaining an access token from DigiKey's API
func (d *DigiKey) authenticate() (string, error) {
	api_url := fmt.Sprintf("%s/v1/oauth2/token", d.BaseURL)

	// Create the application/x-www-form-urlencoded body
	formData := url.Values{}
	formData.Set("client_id", d.ClientID)
	formData.Set("client_secret", d.ClientSecret)
	formData.Set("redirect_uri", "https://localhost")
	formData.Set("grant_type", "client_credentials")

	// Make the POST request
	resp, err := http.Post(api_url, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DigiKey API authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("no access token received from DigiKey API")
	}

	return result.AccessToken, nil
}
