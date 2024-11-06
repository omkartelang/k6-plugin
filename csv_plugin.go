package k6_plugin

import (
	"encoding/csv"
	"go.k6.io/k6/js/modules"
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
