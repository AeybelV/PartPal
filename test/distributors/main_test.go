package distributors

import (
	"testing"

	"github.com/AeybelV/PartPal/internal/utils"
)

var config *utils.Config

func TestMain(m *testing.M) {
	var err error

	// Load configuration from file
	config, err = utils.LoadConfig("../../.partpal.json")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Run the tests
	m.Run()
}
