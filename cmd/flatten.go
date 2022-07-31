package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/cmd/logger"
	"github.com/tuan78/jsonconv/cmd/repository"
)

func NewFlattenCmd() *cobra.Command {
	var (
		lvl int
		gap string
		sm  bool
		sa  bool
	)

	cmd := &cobra.Command{
		Use:   "flatten",
		Short: "Flatten JSON object and JSON array",
		Long:  "Flatten JSON object and JSON array",
		RunE: func(cmd *cobra.Command, args []string) error {
			in := &flattenCmdInput{
				inputPath:  rootFlags.InputPath,
				outputPath: rootFlags.OutputPath,
				raw:        rootFlags.RawData,
				flattenOpt: &jsonconv.FlattenOption{
					Level:     lvl,
					Gap:       gap,
					SkipMap:   sm,
					SkipArray: sa,
				},
			}
			logger := logger.NewLogger(cmd)
			repo := repository.NewRepository()
			return processFlattenCmd(logger, repo, in)
		},
	}

	cmd.PersistentFlags().IntVar(&lvl, "lv", jsonconv.DefaultFlattenLevel, "level for flattening a nested JSON (-1: unlimited, 0: no nested, [1...n]: n level of nested JSON)")
	cmd.PersistentFlags().StringVar(&gap, "ga", jsonconv.DefaultFlattenGap, "gap for separating JSON object with its nested data")
	cmd.PersistentFlags().BoolVar(&sm, "sm", false, "set it true to skip map type")
	cmd.PersistentFlags().BoolVar(&sa, "sa", false, "set it true to skip array type")
	return cmd
}

type flattenCmdInput struct {
	inputPath  string
	outputPath string
	raw        string
	flattenOpt *jsonconv.FlattenOption
}

func processFlattenCmd(logger logger.Logger, repo repository.Repository, in *flattenCmdInput) error {
	var err error

	// Create JSON reader.
	var jr *jsonconv.JsonReader
	switch {
	case in.raw != "":
		jr = jsonconv.NewJsonReader(strings.NewReader(in.raw))
	case in.inputPath != "":
		fi, err := repo.GetFileReader(in.inputPath)
		if err != nil {
			return err
		}
		defer fi.Close()
		jr = jsonconv.NewJsonReader(fi)
	case !repo.IsStdinEmpty():
		fi := repo.GetStdinReader()
		defer fi.Close()
		jr = jsonconv.NewJsonReader(fi)
	default:
		return fmt.Errorf("need to input either raw data, input file path or data from stdin")
	}

	// Read and parse JSON data.
	var encoded interface{}
	err = jr.Read(&encoded)
	if err != nil {
		return fmt.Errorf("invalid JSON data, %v", err)
	}

	var flattened interface{}
	switch val := encoded.(type) {
	case []interface{}:
		var arr jsonconv.JsonArray
		for _, v := range val {
			if obj, ok := v.(jsonconv.JsonObject); ok {
				arr = append(arr, obj)
				continue
			}
			return fmt.Errorf("unsupport type of JSON data")
		}

		// Flatten JSON array.
		for _, obj := range arr {
			jsonconv.FlattenJsonObject(obj, in.flattenOpt)
		}
		flattened = arr
		return outputJsonContent(logger, repo, arr, in.outputPath)
	case jsonconv.JsonObject:
		// Flatten JSON object.
		jsonconv.FlattenJsonObject(val, in.flattenOpt)
		flattened = val
	}

	// Output the JSON content.
	return outputJsonContent(logger, repo, flattened, in.outputPath)
}

func outputJsonContent(logger logger.Logger, repo repository.Repository, data interface{}, filePath string) error {
	// Check and override outputPath if necessary.
	if filePath == "" {
		// Create JSON writer with byte buffer.
		buf := &bytes.Buffer{}
		jw := jsonconv.NewJsonWriter(buf)

		// Write to JSON file.
		err := jw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("%s\n", buf.String())
	} else {
		// Create JSON writer with output file.
		fi, err := repo.CreateFileWriter(filePath)
		if err != nil {
			return err
		}
		defer fi.Close()
		jw := jsonconv.NewJsonWriter(fi)

		// Write to JSON file.
		err = jw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("The JSON file is located at %s\n", filePath)
	}
	return nil
}
