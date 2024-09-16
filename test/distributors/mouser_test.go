package distributors

import (
	"testing"

	"github.com/AeybelV/PartPal/internal/distributors"
)

func TestMouser_Initialize(t *testing.T) {
	mouser := &distributors.Mouser{}
	err := mouser.Initialize(config.Distributors.Mouser.APIKey)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}

	if mouser.APIKey != config.Distributors.Mouser.APIKey {
		t.Errorf("Expected API key to be set correctly")
	}

	err = mouser.Initialize("")
	if err == nil {
		t.Errorf("Expected error when initializing with empty API key, but got none")
	}
}

func TestMouser_QueryMouserPartNumber(t *testing.T) {
	mouser := &distributors.Mouser{}
	err := mouser.Initialize(config.Distributors.Mouser.APIKey)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}
	result, err := mouser.QueryPartNumber("511-STM32C031G4U6")
	if err != nil {
		t.Errorf("Ran into an error with Mouser, got %v", err)
	}

	if result.Description != "ARM Microcontrollers - MCU Mainstream Arm Cortex-M0+ MCU 16 Kbytes Flash 12 Kbytes RAM 48 MHz CPU 2x USART" {
		t.Errorf("The queried products description does not match")
	}
}

func TestMouser_QueryMfrPartNumber(t *testing.T) {
	mouser := &distributors.Mouser{}
	err := mouser.Initialize(config.Distributors.Mouser.APIKey)
	if err != nil {
		t.Errorf("Recieved error when initializing, got %v", err)
	}

	result, err := mouser.QueryPartNumber("STM32C031G4U6")
	if err != nil {
		t.Errorf("Ran into an error with Mouser, got %v", err)
	}

	if result.Description != "ARM Microcontrollers - MCU Mainstream Arm Cortex-M0+ MCU 16 Kbytes Flash 12 Kbytes RAM 48 MHz CPU 2x USART" {
		t.Errorf("The queried products description does not match")
	}
}
