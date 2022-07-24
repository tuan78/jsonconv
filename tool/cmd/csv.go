package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/tool/params"
	"github.com/tuan78/jsonconv/utils"
)

func NewCsvCmd() *cobra.Command {
	var (
		baseHeaders []string
		delimiter   string
		useCRLF     bool
	)

	cmd := &cobra.Command{
		Use:   "csv",
		Short: "Convert JSON to CSV",
		Long:  "Convert JSON to CSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processCsvCmd(baseHeaders, delimiter, useCRLF)
		},
	}

	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().StringSliceVar(&baseHeaders, "headers", nil, "headers in CSV that always appears before dynamic headers (auto detected from JSON)")
	cmd.PersistentFlags().StringVar(&delimiter, "delimiter", ",", "field delimiter")
	cmd.PersistentFlags().BoolVar(&useCRLF, "useCRLF", false, "set it true to use \\r\\n as the line terminator")
	return cmd
}

func processCsvCmd(baseHeaders []string, delimiter string, useCRLF bool) error {
	var err error

	// Create JSON reader.
	var jsonReader *jsonconv.JsonReader
	switch {
	case params.RawData != "":
		jsonReader = jsonconv.NewJsonReaderFromString(params.RawData)
	case params.InputPath != "":
		jsonReader, err = jsonconv.NewJsonReaderFromFile(params.InputPath)
		if err != nil {
			return err
		}
	case !utils.IsStdinEmpty():
		jsonReader = jsonconv.NewJsonReader(os.Stdin)
	default:
		return fmt.Errorf("need to input either raw data, input file path or data from stdin")
	}

	// Read JSON data and store in jsonArray.
	fmt.Println("Processing...")
	var encoded interface{}
	err = jsonReader.Read(&encoded)
	if err != nil {
		return err
	}
	var jsonArray jsonconv.JsonArray
	switch val := encoded.(type) {
	case []interface{}:
		for _, v := range val {
			if jsonObject, ok := v.(jsonconv.JsonObject); ok {
				jsonArray = append(jsonArray, jsonObject)
				continue
			}
			return fmt.Errorf("unknown type of JSON data")
		}
	case jsonconv.JsonObject:
		jsonArray = append(jsonArray, val)
	default:
		return fmt.Errorf("unknown type of JSON data")
	}

	// Convert JSON to CSV.
	var csvData [][]string
	csvData, err = jsonconv.Convert(&jsonconv.ConvertInput{
		JsonArray:    jsonArray,
		FlattenLevel: params.FlattenLevel,
		BaseHeaders:  baseHeaders,
	})
	if err != nil {
		return err
	}
	if len(csvData) == 0 {
		return fmt.Errorf("empty CSV data")
	}

	return writeToCsvFile(csvData, delimiter, useCRLF)
}

func writeToCsvFile(csvData jsonconv.CsvData, delimiter string, useCRLF bool) error {
	// Check and override outputPath if necessary.
	if params.OutputPath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("output-%v.csv", time.Now().UTC().Unix())
		params.OutputPath = filepath.Join(workingDir, fileName)
	}

	// Create CSV writer.
	runes := []rune(delimiter)
	var delimiterRune *rune
	if len(runes) > 0 {
		delimiterRune = &runes[0]
	}
	csvWriter, err := jsonconv.NewCsvWriterFromFile(params.OutputPath)
	if err != nil {
		return err
	}
	csvWriter.Delimiter = delimiterRune
	csvWriter.UseCRLF = useCRLF

	// Write to CSV file.
	err = csvWriter.Write(csvData)
	if err != nil {
		return err
	}

	fmt.Printf("Done. The CSV file is located at %s\n", params.OutputPath)
	return nil
}
