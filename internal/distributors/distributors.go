package distributors

// PartInfo represents the unified information returned from querying a part number.
type PartInfo struct {
	PartNumber             string  `json:"part_number"`
	ManufacturerPartNumber string  `json:"mfr_part_number"`
	Manufacturer           string  `json:"manufacturer"`
	Description            string  `json:"description"`
	UnitPrice              float64 `json:"unitprice"`
	Availability           int     `json:"availability"`
	Quantity               int
	ProductURL             string `json:"product_url"`
	DataSheetURL           string `json:"datasheet_url"`
}
type Distributor interface {
	Initialize(params ...string) error
	QueryPartNumber(partNumber string) (PartInfo, error)
}
