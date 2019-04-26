package main

// Transforms Refcardz CSV to Refcardz JSON (dzone-refcardz-direct-dl.csv --> dzone-refcardz-direct-dl.json)
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// Refcard struct contains refcard numbers, file names and direct URLs
type Refcard struct {
	Number int
	Name   string
	URL    string
}

func main() {
	csvFile, err := os.Open("./input.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var ref Refcard
	var refcardz []Refcard

	for i, each := range csvData {
		i++
		ref.Number = i
		ref.Name = each[0]
		ref.URL = each[1]
		refcardz = append(refcardz, ref)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(refcardz)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))

	jsonFile, err := os.Create("./output.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
