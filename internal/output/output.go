package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Format represents an output format
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatPlain Format = "plain"
)

// Formatter handles output formatting
type Formatter struct {
	format Format
	writer io.Writer
	color  bool
}

// New creates a new formatter
func New(format Format, writer io.Writer) *Formatter {
	return &Formatter{
		format: format,
		writer: writer,
		color:  isTerminal(writer),
	}
}

// isTerminal checks if the writer is a terminal
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fileInfo, err := f.Stat()
		if err != nil {
			return false
		}
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// Output writes data in the configured format
func (f *Formatter) Output(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.outputJSON(data)
	case FormatCSV:
		return f.outputCSV(data)
	case FormatTable:
		return f.outputTable(data)
	case FormatPlain:
		return f.outputPlain(data)
	default:
		return fmt.Errorf("unsupported format: %s", f.format)
	}
}

// outputJSON outputs data as JSON
func (f *Formatter) outputJSON(data interface{}) error {
	enc := json.NewEncoder(f.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// outputCSV outputs data as CSV
func (f *Formatter) outputCSV(data interface{}) error {
	w := csv.NewWriter(f.writer)
	defer w.Flush()

	// Convert data to slice of maps
	rows, headers := f.dataToRows(data)
	if len(rows) == 0 {
		return nil
	}

	// Write headers
	if err := w.Write(headers); err != nil {
		return err
	}

	// Write rows
	for _, row := range rows {
		record := make([]string, len(headers))
		for i, header := range headers {
			record[i] = fmt.Sprint(row[header])
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// outputTable outputs data as an aligned table
func (f *Formatter) outputTable(data interface{}) error {
	rows, headers := f.dataToRows(data)
	if len(rows) == 0 {
		return nil
	}

	table := tablewriter.NewWriter(f.writer)
	table.Header(headers)

	for _, row := range rows {
		record := make([]string, len(headers))
		for i, header := range headers {
			val := fmt.Sprint(row[header])
			if f.color && header == "PRIORITY" {
				val = f.colorPriority(val)
			}
			record[i] = val
		}
		table.Append(record)
	}

	table.Render()
	return nil
}

// outputPlain outputs data as plain text (one value per line)
func (f *Formatter) outputPlain(data interface{}) error {
	rows, _ := f.dataToRows(data)
	for _, row := range rows {
		// Output first field value only
		for _, val := range row {
			fmt.Fprintln(f.writer, val)
			break
		}
	}
	return nil
}

// dataToRows converts data to rows and headers
func (f *Formatter) dataToRows(data interface{}) ([]map[string]interface{}, []string) {
	// Handle slice of structs or maps
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		// Single item, wrap in slice
		v = reflect.ValueOf([]interface{}{data})
	}

	if v.Len() == 0 {
		return nil, nil
	}

	rows := make([]map[string]interface{}, 0, v.Len())
	headers := []string{}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		row := f.structToMap(item)
		rows = append(rows, row)

		// Collect headers from first item
		if i == 0 {
			for k := range row {
				headers = append(headers, k)
			}
		}
	}

	return rows, headers
}

// structToMap converts a struct to a map
func (f *Formatter) structToMap(item interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Try JSON marshaling as a simple approach
	jsonData, err := json.Marshal(item)
	if err != nil {
		return result
	}

	var m map[string]interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return result
	}

	// Flatten nested objects for display
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			// For nested objects, just use a representative field
			if name, ok := vm["name"]; ok {
				result[strings.ToUpper(k)] = name
			} else if id, ok := vm["id"]; ok {
				result[strings.ToUpper(k)] = id
			} else {
				result[strings.ToUpper(k)] = "<object>"
			}
		} else if v != nil {
			result[strings.ToUpper(k)] = v
		}
	}

	return result
}

// colorPriority colors priority values
func (f *Formatter) colorPriority(val string) string {
	if !f.color {
		return val
	}

	switch {
	case strings.Contains(val, "P0") || strings.Contains(val, "Urgent"):
		return color.RedString(val)
	case strings.Contains(val, "P1") || strings.Contains(val, "High"):
		return color.YellowString(val)
	case strings.Contains(val, "P2") || strings.Contains(val, "Medium"):
		return color.CyanString(val)
	default:
		return val
	}
}
