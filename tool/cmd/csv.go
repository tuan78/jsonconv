package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuan78/jsonconv"
	"github.com/tuan78/jsonconv/tool/params"
	"github.com/tuan78/jsonconv/utils"
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
				inputPath:  params.InputPath,
				outputPath: params.OutputPath,
				raw:        params.RawData,
				baseHs:     baseHs,
				delim:      delim,
				useCRLF:    crlf,
			}
			if !noft {
				in.flattenOp = &jsonconv.FlattenOption{
					Level:     flv,
					Gap:       fga,
					SkipMap:   fsm,
					SkipArray: fsa,
				}
			}
			return processCsvCmd(in)
		},
	}

	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().StringSliceVar(&baseHs, "hs", nil, "headers in CSV that always appears before dynamic headers (auto detected from JSON)")
	cmd.PersistentFlags().StringVar(&delim, "delim", ",", "field delimiter")
	cmd.PersistentFlags().BoolVar(&crlf, "crlf", false, "set it true to use \\r\\n as the line terminator")
	cmd.PersistentFlags().BoolVar(&noft, "noft", false, "set it true to skip JSON flattening")
	cmd.PersistentFlags().IntVar(&flv, "flv", jsonconv.FlattenLevelDefault, "flatten level for flattening a nested JSON (-1: unlimited, 0: no nested, [1...n]: n level of nested JSON)")
	cmd.PersistentFlags().StringVar(&fga, "fga", jsonconv.FlattenGapDefault, "flatten gap for separating JSON object with its nested data")
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
	flattenOp  *jsonconv.FlattenOption
}

func processCsvCmd(in *csvCmdInput) error {
	var err error

	// Create JSON reader.
	var jr *jsonconv.JsonReader
	switch {
	case in.raw != "":
		jr = jsonconv.NewJsonReader(strings.NewReader(in.raw))
	case in.inputPath != "":
		fi, err := os.Open(in.inputPath)
		if err != nil {
			return err
		}
		defer fi.Close()
		jr = jsonconv.NewJsonReader(fi)
	case !utils.IsStdinEmpty():
		fi := os.Stdin
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
	var arr jsonconv.JsonArray
	switch val := encoded.(type) {
	case []interface{}:
		for _, v := range val {
			if obj, ok := v.(jsonconv.JsonObject); ok {
				arr = append(arr, obj)
				continue
			}
			return fmt.Errorf("unknown type of JSON data")
		}
	case jsonconv.JsonObject:
		arr = append(arr, val)
	default:
		return fmt.Errorf("unknown type of JSON data")
	}

	// Convert JSON to CSV.
	data := jsonconv.ToCsv(arr, &jsonconv.ToCsvOption{
		FlattenOption: in.flattenOp,
		BaseHeaders:   in.baseHs,
	})
	if len(data) == 0 {
		return fmt.Errorf("empty CSV data")
	}

	// Convert in.delim to rune.
	runes := []rune(in.delim)
	var delimRune *rune
	if len(runes) > 0 {
		delimRune = &runes[0]
	}

	// Output the CSV content.
	return outputCsvContent(data, in.outputPath, delimRune, in.useCRLF)
}

func outputCsvContent(data jsonconv.CsvData, filePath string, delim *rune, useCRLF bool) error {
	var err error

	// Check and override outputPath if necessary.
	path := filePath
	if path == "" {
		// Create CSV writer with byte buffer.
		buf := &bytes.Buffer{}
		cw := jsonconv.NewCsvWriter(buf)
		cw.Delimiter = delim
		cw.UseCRLF = useCRLF

		// Write to CSV file.
		err = cw.Write(data)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", buf.String())
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

		// Create CSV writer with output file.
		cw := jsonconv.NewCsvWriter(fi)
		cw.Delimiter = delim
		cw.UseCRLF = useCRLF

		// Write to CSV file.
		err = cw.Write(data)
		if err != nil {
			return err
		}
		fmt.Printf("The CSV file is located at %s\n", path)
	}
	return nil
}
