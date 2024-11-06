package k6_plugin

import (
	"bufio"
	"encoding/csv"
	"go.k6.io/k6/js/modules"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

// Register the plugin as a k6 module
func init() {
	modules.Register("k6/x/k6plugin", new(K6Plugin))
}

// K6Plugin is the main structure for the plugin
type K6Plugin struct{}

// New creates a new instance of the csv and returns it.
func (p *K6Plugin) CreateCSVWriter(csvFile string) *csv.Writer {

	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	writercsv := csv.NewWriter(file)
	defer writercsv.Flush()
	return writercsv
}

// WriteCSVHeader writes the response to the specified CSV writer as a header.
func (p *K6Plugin) WriteCSVHeader(writer *csv.Writer, response string) {
	p.writeCSV(writer, response)
}

// WriteResponse writes the response to the specified CSV writer.
func (p *K6Plugin) WriteResponse(writer *csv.Writer, response string) {
	p.writeCSV(writer, response)
}

// writeCSV is a helper function to write a response to a CSV writer.
func (p *K6Plugin) writeCSV(writer *csv.Writer, response string) {
	row := strings.Split(response, ",")
	writer.Write(row)
}

// IsFileEmpty checks if the specified file is empty.
func (p *K6Plugin) IsFileEmpty(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return fileInfo.Size() == 0, nil
}

// CalculatePercentiles calculates the 50th, 90th, 95th, and 99th percentiles of the data.
func (p *K6Plugin) CalculatePercentiles(data []float64, count int) map[int]float64 {
	percentiles := map[int]float64{50: 0, 90: 0, 95: 0, 99: 0}
	if count == 0 {
		return percentiles
	}

	sort.Float64s(data)

	percentiles[50] = data[int(float64(count)*0.50)]
	percentiles[90] = data[int(float64(count)*0.90)-1]
	percentiles[95] = data[int(float64(count)*0.95)-1]
	percentiles[99] = data[int(float64(count)*0.99)-1]

	return percentiles
}

// CalculateStdDev calculates the standard deviation of the data.
func (p *K6Plugin) CalculateStdDev(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0.0
	}

	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(n)

	variance := 0.0
	for _, v := range data {
		variance += math.Pow(v-mean, 2)
	}
	variance /= float64(n)

	return math.Sqrt(variance)
}

// WriteString writes string to file
func (p *K6Plugin) WriteString(path string, s string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(s); err != nil {
		return err
	}
	return nil
}

// AppendString appends string to file
func (p *K6Plugin) AppendString(path string, s string) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(s); err != nil {
		return err
	}
	return nil
}

// WriteBytes writes binary file
func (p *K6Plugin) WriteBytes(path string, b []byte) error {
	err := os.WriteFile(path, b, 0o644)
	if err != nil {
		return err
	}
	return nil
}

// ClearFile removes all the contents of a file
func (p *K6Plugin) ClearFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}
	return nil
}

// RenameFile renames file from oldPath to newPath
func (p K6Plugin) RenameFile(oldPath string, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes file
func (p *K6Plugin) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// RemoveRowsBetweenValues removes the rows from a file between start and end (inclusive)
func (p *K6Plugin) RemoveRowsBetweenValues(path string, start, end int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	lineCount := 0

	// Read the entire file into memory
	for scanner.Scan() {
		lineCount++
		if lineCount < start || lineCount > end {
			lines = append(lines, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write the modified contents back to the file
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
