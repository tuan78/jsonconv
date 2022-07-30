package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/tool/logger"
	"github.com/tuan78/jsonconv/tool/params"
	"github.com/tuan78/jsonconv/tool/repository"
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
				inputPath:  params.InputPath,
				outputPath: params.OutputPath,
				raw:        params.RawData,
				flattenOpt: &jsonconv.FlattenOption{
					Level:     lvl,
					Gap:       gap,
					SkipMap:   sm,
					SkipArray: sa,
				},
			}
			logger := logger.NewCmdLogger(cmd)
			repo := repository.NewRepository()
			return processFlattenCmd(logger, repo, in)
		},
	}

	cmd.PersistentFlags().IntVar(&lvl, "lv", jsonconv.DefaultFlattenLevel, "level for flattening a nested JSON (-1: unlimited, 0: no nested, [1...n]: n level of nested JSON)")
	cmd.PersistentFlags().StringVar(&gap, "ga", jsonconv.DefaultFlattenGap, "gap for separating JSON object with its nested data")
	cmd.PersistentFlags().BoolVar(&sm, "sm", false, "skip map type")
	cmd.PersistentFlags().BoolVar(&sa, "sa", false, "skip array type")
	return cmd
}

type flattenCmdInput struct {
	inputPath  string
	outputPath string
	raw        string
	flattenOpt *jsonconv.FlattenOption
}

func processFlattenCmd(logger logger.CmdLogger, repo repository.Repository, in *flattenCmdInput) error {
	var err error

	// Create JSON reader.
	var jr *jsonconv.JsonReader
	switch {
	case in.raw != "":
		jr = jsonconv.NewJsonReader(strings.NewReader(in.raw))
	case in.inputPath != "":
		fi, err := repo.GetFileReadCloser(in.inputPath)
		if err != nil {
			return err
		}
		defer fi.Close()
		jr = jsonconv.NewJsonReader(fi)
	case !repo.IsStdinEmpty():
		fi := repo.GetStdinReadCloser()
		defer fi.Close()
		jr = jsonconv.NewJsonReader(fi)
	default:
		return fmt.Errorf("need to input either raw data, input file path or data from stdin")
	}

	// Read and parse JSON data.
	var encoded interface{}
	err = jr.Read(&encoded)
	if err != nil {
		return err
	}

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

		// Output the JSON content.
		return outputJsonContent(logger, arr, in.outputPath)
	case jsonconv.JsonObject:
		// Flatten JSON object.
		jsonconv.FlattenJsonObject(val, in.flattenOpt)

		// Output the JSON content.
		return outputJsonContent(logger, val, in.outputPath)
	}

	return nil
}

func outputJsonContent(logger logger.CmdLogger, data interface{}, filePath string) error {
	var err error

	// Check and override outputPath if necessary.
	path := filePath
	if path == "" {
		// Create JSON writer with byte buffer.
		buf := &bytes.Buffer{}
		jw := jsonconv.NewJsonWriter(buf)

		// Write to JSON file.
		err = jw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("%s\n", buf.String())
	} else {
		var fi *os.File
		// Check file path and make dir accordingly.
		if strings.Contains(path, string(filepath.Separator)) ||
			strings.HasPrefix(path, ".") ||
			strings.HasPrefix(path, "~") {
			// Ensure all dir in path exists.
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}
			fi, err = os.Create(path)
			if err != nil {
				return err
			}
			defer fi.Close()
		} else {
			// Path is only file name so override it with full path (working dir + file name).
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			path = filepath.Join(dir, filePath)
			fi, err = os.Create(path)
			if err != nil {
				return err
			}
			defer fi.Close()
		}

		// Create JSON writer with output file.
		jw := jsonconv.NewJsonWriter(fi)

		// Write to JSON file.
		err = jw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("The JSON file is located at %s\n", path)
	}
	return nil
}
