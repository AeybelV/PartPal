package distributors

import (
	"testing"

	"github.com/AeybelV/PartPal/internal/distributors"
)

func TestDigiKey_Initialize(t *testing.T) {
	digikey := &distributors.DigiKey{}
	err := digikey.Initialize(config.Distributors.DigiKey.ClientID, config.Distributors.DigiKey.ClientSecret)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}

	if digikey.ClientID != config.Distributors.DigiKey.ClientID {
		t.Errorf("Expected ClientID to be set correctly")
	}

	if digikey.ClientSecret != config.Distributors.DigiKey.ClientSecret {
		t.Errorf("Expected Client Secret to be set correctly")
	}

	err = digikey.Initialize("")
	if err == nil {
		t.Errorf("Expected error when initializing with empty API key, but got none")
	}
}

func TestDigiKey_QuerydigikeyPartNumber(t *testing.T) {
	digikey := &distributors.DigiKey{}
	err := digikey.Initialize(config.Distributors.DigiKey.ClientID, config.Distributors.DigiKey.ClientSecret)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}
	result, err := digikey.QueryPartNumber("497-STM32C031G4U6-ND")
	if err != nil {
		t.Errorf("Ran into an error with digikey, got %v", err)
	}

	if result.Description != "IC MCU 32BIT 16KB FLASH 28UFQFPN" {
		t.Errorf("The queried products description does not match")
	}
}

func TestDigiKey_QueryMfrPartNumber(t *testing.T) {
	digikey := &distributors.DigiKey{}
	err := digikey.Initialize(config.Distributors.DigiKey.ClientID, config.Distributors.DigiKey.ClientSecret)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}

	result, err := digikey.QueryPartNumber("STM32C031G4U6")
	if err != nil {
		t.Errorf("Ran into an error with digikey, got %v", err)
	}

	if result.Description != "IC MCU 32BIT 16KB FLASH 28UFQFPN" {
		t.Errorf("The queried products description does not match")
	}
}
