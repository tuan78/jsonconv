package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/internal/cli/logger"
	"github.com/tuan78/jsonconv/internal/cli/repository"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		data, err := io.ReadAll(fi)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		jr = jsonconv.NewJsonReader(bytes.NewReader(data))
	case !repo.IsStdinEmpty():
		fi := repo.GetStdinReader()
		defer fi.Close()
		data, err := io.ReadAll(fi)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		jr = jsonconv.NewJsonReader(bytes.NewReader(data))
	default:
		return fmt.Errorf("need to input either raw data, input file path or data from stdin")
	}

	// Read and parse JSON data.
	var encoded any
	err = jr.Read(&encoded)
	if err != nil {
		return fmt.Errorf("invalid JSON data, %v", err)
	}

	switch val := encoded.(type) {
	case []any:
		var arr []map[string]any
		for _, v := range val {
			if obj, ok := v.(map[string]any); ok {
				arr = append(arr, obj)
				continue
			}
			return fmt.Errorf("unsupport type of JSON data")
		}

		// Flatten JSON array.
		for _, obj := range arr {
			jsonconv.Flatten(obj, in.flattenOpt)
		}
		return outputJsonContent(logger, repo, arr, in.outputPath)
	case map[string]any:
		// Flatten JSON object.
		jsonconv.Flatten(val, in.flattenOpt)
		return outputJsonContent(logger, repo, val, in.outputPath)
	}

	return fmt.Errorf("unsupported JSON data type")
}

func outputJsonContent(logger logger.Logger, repo repository.Repository, data any, filePath string) error {
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
