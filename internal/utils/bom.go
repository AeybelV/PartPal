package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/AeybelV/PartPal/internal/distributors"
)

// A Struct representing a ROW in a BOM thats provided
type BOMRow struct {
	PartNumber  string
	Description string
	Quantity    int
}

// This map will define alternative field names for each column we're interested in
var BOMFieldMapping = map[string]string{
	"Part Number":    "PartNumber",
	"Component":      "PartNumber",
	"PN":             "PartNumber",
	"P/N":            "PartNumber",
	"Product Number": "PartNumber",
	"Part":           "PartNumber",
	"Description":    "Description",
	"Desc":           "Description",
	"Desc.":          "Description",
	"Qty":            "Quantity",
	"Qty.":           "Quantity",
	"Quantity":       "Quantity",
}

// Struct that PartPal works with internally to represent BOMs
type PartPalBOM struct {
	Components []distributors.PartInfo
	TotalCost  float64
}

// QueryResult struct to hold the result of querying a distributor
type QueryResult struct {
	PartInfo distributors.PartInfo
	Err      error
}

// ExportBOMToCSV takes a BOM and exports it to a CSV file
func ExportBOMToCSV(bom PartPalBOM, filepath string) error {
	// Create the CSV file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Dynamically get the header from PartInfo struct field names
	partInfoType := reflect.TypeOf(distributors.PartInfo{})
	header := make([]string, partInfoType.NumField())
	for i := 0; i < partInfoType.NumField(); i++ {
		header[i] = partInfoType.Field(i).Name
	}

	// Write the header to the CSV
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each part in the BOM
	for _, part := range bom.Components {
		// Dynamically get the field values for each part
		row := make([]string, partInfoType.NumField())
		partValue := reflect.ValueOf(part)
		for i := 0; i < partInfoType.NumField(); i++ {
			field := partValue.Field(i)
			switch field.Kind() {
			case reflect.String:
				row[i] = field.String()
			case reflect.Float64:
				row[i] = fmt.Sprintf("%.2f", field.Float())
			case reflect.Int:
				row[i] = strconv.Itoa(int(field.Int()))
			default:
				row[i] = fmt.Sprintf("%v", field.Interface())
			}
		}

		// Write the row to the CSV file
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// ReadBOM function to read from a CSV file and return a list of BOMRows
func ReadBOMCSV(filepath string) ([]BOMRow, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header row to map column names
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create a map of column indices to unified field names
	columnMap := make(map[int]string)
	for i, column := range header {
		normalized := strings.TrimSpace(column)
		if unifiedField, ok := BOMFieldMapping[normalized]; ok {
			columnMap[i] = unifiedField
		}
	}

	var bomRows []BOMRow

	// Process each row and map the fields
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		row := BOMRow{}
		// Use reflection to set struct fields dynamically
		rowValue := reflect.ValueOf(&row).Elem()

		for i, value := range record {
			if fieldName, ok := columnMap[i]; ok {
				// Find the field by name
				field := rowValue.FieldByName(fieldName)
				if !field.IsValid() {
					continue
				}

				// Set the field value depending on its type
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Int:
					intVal, err := strconv.Atoi(value)
					if err != nil {
						return nil, fmt.Errorf("failed to convert quantity to int: %w", err)
					}
					field.SetInt(int64(intVal))
				}
			}
		}
		bomRows = append(bomRows, row)
	}

	return bomRows, nil
}

// Wrapper function to read BOM and convert it to PartInfo array
func ReadBOM(filepath, filetype string) (PartPalBOM, error) {
	var bomRows []BOMRow
	var err error
	ppBOM := PartPalBOM{}
	// Determine which ReadBOM function to call based on file type
	switch strings.ToLower(filetype) {
	case "csv":
		bomRows, err = ReadBOMCSV(filepath) // Calls the CSV reader function
		if err != nil {
			return ppBOM, fmt.Errorf("failed to read BOM: %w", err)
		}
	default:
		return ppBOM, fmt.Errorf("unsupported file type: %s", filetype)
	}

	// Convert BOMRow array to PartInfo array
	partInfos := make([]distributors.PartInfo, len(bomRows))
	for i, bom := range bomRows {
		qty := 1
		if bom.Quantity != 0 {
			qty = bom.Quantity
		}
		partInfos[i] = distributors.PartInfo{
			PartNumber:  bom.PartNumber,
			Description: bom.Description,
			Quantity:    qty,
			// Add more mappings if needed (e.g., from BOMRow to PartInfo)
		}
	}

	ppBOM.Components = partInfos
	ppBOM.TotalCost = 0.0

	return ppBOM, nil
}

// OptimizeBOM function that takes a list of parts and distributors and returns the optimized BOM
func OptimizeBOM_Serial(bom PartPalBOM, distributorsList []distributors.Distributor) (PartPalBOM, error) {
	var optimizedParts []distributors.PartInfo
	totalCost := 0.0

	// Iterate over each part in the BOM
	for _, part := range bom.Components {
		var bestPartInfo distributors.PartInfo
		lowestCost := math.MaxFloat64 // Start with an arbitrarily high price

		// Query each distributor for the part
		for _, distributor := range distributorsList {
			partInfo, err := distributor.QueryPartNumber(part.PartNumber)
			if err != nil {
				// If distributor can't supply the part, skip
				// fmt.Printf("Cant provide part")
				continue
			}

			// Check if the distributor offers a lower price
			if partInfo.UnitPrice < lowestCost {
				lowestCost = partInfo.UnitPrice
				bestPartInfo = partInfo
				bestPartInfo.Quantity = part.Quantity
			}
		}

		// If we found a cheaper part, update the part info
		if bestPartInfo.PartNumber != "" {
			optimizedParts = append(optimizedParts, bestPartInfo)
			totalCost += bestPartInfo.UnitPrice * float64(part.Quantity)
		} else {
			// If no distributor could supply the part, retain the original part info
			optimizedParts = append(optimizedParts, part)
			totalCost += part.UnitPrice * float64(part.Quantity)
		}

	}

	// Create an OptimizedBOM struct to hold the results
	optimizedBOM := PartPalBOM{
		Components: optimizedParts,
		TotalCost:  totalCost,
	}

	return optimizedBOM, nil
}

// OptimizeBOM function that takes a list of parts and distributors and returns the optimized BOM
func OptimizeBOM_Parallel(bom PartPalBOM, distributorsList []distributors.Distributor) (PartPalBOM, error) {
	var optimizedParts []distributors.PartInfo
	var totalCost float64
	var mu sync.Mutex // Mutex to protect shared data

	var wg sync.WaitGroup
	optimizedPartsChan := make(chan distributors.PartInfo, len(bom.Components))

	// Iterate over each part in the BOM in parallel
	for _, part := range bom.Components {
		wg.Add(1)

		go func(part distributors.PartInfo) {
			defer wg.Done()

			var bestPartInfo distributors.PartInfo
			lowestCost := math.MaxFloat64 // Start with an arbitrarily high price

			var partWG sync.WaitGroup
			partResults := make(chan QueryResult, len(distributorsList))

			// Query each distributor in parallel for the part
			for _, distributor := range distributorsList {
				partWG.Add(1)
				go func(d distributors.Distributor) {
					defer partWG.Done()
					partInfo, err := d.QueryPartNumber(part.PartNumber)
					partResults <- QueryResult{PartInfo: partInfo, Err: err}
				}(distributor)
			}

			// Close the partResults channel when all queries for the part are done
			go func() {
				partWG.Wait()
				close(partResults)
			}()

			// Collect the results for the part and find the lowest price
			for result := range partResults {
				if result.Err == nil && result.PartInfo.UnitPrice < lowestCost {
					lowestCost = result.PartInfo.UnitPrice
					bestPartInfo = result.PartInfo
					bestPartInfo.Quantity = part.Quantity
				}
			}

			// Lock to safely update shared state (total cost and parts)
			mu.Lock()
			defer mu.Unlock()

			// If we found a cheaper part, update the part info
			if bestPartInfo.PartNumber != "" {
				optimizedPartsChan <- bestPartInfo
			} else {
				// If no distributor could supply the part, retain the original part info
				optimizedPartsChan <- part
			}
		}(part)
	}

	// Close the channel when all parts have been processed
	go func() {
		wg.Wait()
		close(optimizedPartsChan)
	}()

	// Collect all optimized parts from the channel
	for part := range optimizedPartsChan {
		optimizedParts = append(optimizedParts, part)
		totalCost += part.UnitPrice * float64(part.Quantity)
	}

	// Create an OptimizedBOM struct to hold the results
	optimizedBOM := PartPalBOM{
		Components: optimizedParts,
		TotalCost:  totalCost,
	}

	return optimizedBOM, nil
}

// OptimizeBOM function that takes a list of parts and distributors and returns the optimized BOM
func OptimizeBOM(bom PartPalBOM, distributorsList []distributors.Distributor) (PartPalBOM, error) {
	return OptimizeBOM_Serial(bom, distributorsList)
}
