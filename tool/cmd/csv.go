package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/tool/logger"
	"github.com/tuan78/jsonconv/tool/repository"
)

func NewCsvCmd() *cobra.Command {
	var (
		baseHs []string
		delim  string
		crlf   bool
		noft   bool
		flv    int
		fga    string
		fsm    bool
		fsa    bool
	)

	cmd := &cobra.Command{
		Use:   "csv",
		Short: "Convert JSON to CSV",
		Long:  "Convert JSON to CSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			in := &csvCmdInput{
				inputPath:  rootFlags.InputPath,
				outputPath: rootFlags.OutputPath,
				raw:        rootFlags.RawData,
				baseHs:     baseHs,
				delim:      delim,
				useCRLF:    crlf,
			}
			if !noft {
				in.flattenOpt = &jsonconv.FlattenOption{
					Level:     flv,
					Gap:       fga,
					SkipMap:   fsm,
					SkipArray: fsa,
				}
			}
			logger := logger.NewLogger(cmd)
			repo := repository.NewRepository()
			return processCsvCmd(logger, repo, in)
		},
	}

	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().StringSliceVar(&baseHs, "hs", nil, "headers in CSV that always appears before dynamic headers (auto detected from JSON)")
	cmd.PersistentFlags().StringVar(&delim, "delim", ",", "field delimiter")
	cmd.PersistentFlags().BoolVar(&crlf, "crlf", false, "set it true to use \\r\\n as the line terminator")
	cmd.PersistentFlags().BoolVar(&noft, "noft", false, "set it true to skip JSON flattening")
	cmd.PersistentFlags().IntVar(&flv, "flv", jsonconv.DefaultFlattenLevel, "flatten level for flattening a nested JSON (-1: unlimited, 0: no nested, [1...n]: n level of nested JSON)")
	cmd.PersistentFlags().StringVar(&fga, "fga", jsonconv.DefaultFlattenGap, "flatten gap for separating JSON object with its nested data")
	cmd.PersistentFlags().BoolVar(&fsm, "fsm", false, "flatten but skip map type")
	cmd.PersistentFlags().BoolVar(&fsa, "fsa", false, "flatten but skip array type")
	return cmd
}

type csvCmdInput struct {
	inputPath  string
	outputPath string
	raw        string
	baseHs     []string
	delim      string
	useCRLF    bool
	flattenOpt *jsonconv.FlattenOption
}

func processCsvCmd(logger logger.Logger, repo repository.Repository, in *csvCmdInput) error {
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
	var arr jsonconv.JsonArray
	switch val := encoded.(type) {
	case []interface{}:
		for _, v := range val {
			obj, ok := v.(jsonconv.JsonObject)
			if !ok {
				return fmt.Errorf("unsupport type of JSON data")
			}
			if len(obj) == 0 {
				continue
			}
			arr = append(arr, obj)
		}
	case jsonconv.JsonObject:
		if len(val) != 0 {
			arr = append(arr, val)
		}
	}

	// Convert JSON to CSV.
	data := jsonconv.ToCsv(arr, &jsonconv.ToCsvOption{
		FlattenOption: in.flattenOpt,
		BaseHeaders:   in.baseHs,
	})

	// Convert in.delim to rune.
	runes := []rune(in.delim)
	var delimRune *rune
	if len(runes) > 0 {
		delimRune = &runes[0]
	}

	// Output the CSV content.
	return outputCsvContent(logger, repo, data, in.outputPath, delimRune, in.useCRLF)
}

func outputCsvContent(logger logger.Logger, repo repository.Repository, data jsonconv.CsvData, filePath string, delim *rune, useCRLF bool) error {
	// Check and override outputPath if necessary.
	if filePath == "" {
		// Create CSV writer with byte buffer.
		buf := &bytes.Buffer{}
		cw := jsonconv.NewCsvWriter(buf)
		if delim != nil {
			cw.Delimiter = *delim
		}
		cw.UseCRLF = useCRLF

		// Write to CSV file.
		err := cw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("%s\n", buf.String())
	} else {
		// Create CSV writer with output file.
		fi, err := repo.CreateFileWriter(filePath)
		if err != nil {
			return err
		}
		defer fi.Close()
		cw := jsonconv.NewCsvWriter(fi)
		if delim != nil {
			cw.Delimiter = *delim
		}
		cw.UseCRLF = useCRLF

		// Write to CSV file.
		err = cw.Write(data)
		if err != nil {
			return err
		}
		logger.Printf("The CSV file is located at %s\n", filePath)
	}
	return nil
}
