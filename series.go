// MIT License

// Copyright (c) 2022 Tree Xie

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package charts

import (
	"math"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type SeriesData struct {
	// The value of series data
	Value float64
	// The style of series data
	Style chart.Style
}

func NewSeriesFromValues(values []float64, chartType ...string) Series {
	s := Series{
		Data: NewSeriesDataFromValues(values),
	}
	if len(chartType) != 0 {
		s.Type = chartType[0]
	}
	return s
}

func NewSeriesDataFromValues(values []float64) []SeriesData {
	data := make([]SeriesData, len(values))
	for index, value := range values {
		data[index] = SeriesData{
			Value: value,
		}
	}
	return data
}

type SeriesLabel struct {
	// Data label formatter, which supports string template.
	// {b}: the name of a data item.
	// {c}: the value of a data item.
	// {d}: the percent of a data item(pie chart).
	Formatter string
	// The color for label
	Color drawing.Color
	// Show flag for label
	Show bool
	// Distance to the host graphic element.
	Distance int
}

const (
	SeriesMarkDataTypeMax     = "max"
	SeriesMarkDataTypeMin     = "min"
	SeriesMarkDataTypeAverage = "average"
)

type SeriesMarkData struct {
	// The mark data type, it can be "max", "min", "average".
	// The "average" is only for mark line
	Type string
}
type SeriesMarkPoint struct {
	// The width of symbol, default value is 30
	SymbolSize int
	// The mark data of series mark point
	Data []SeriesMarkData
}
type SeriesMarkLine struct {
	// The mark data of series mark line
	Data []SeriesMarkData
}
type Series struct {
	index int
	// The type of series, it can be "line", "bar" or "pie".
	// Default value is "line"
	Type string
	// The data list of series
	Data []SeriesData
	// The Y axis index, it should be 0 or 1.
	// Default value is 1
	YAxisIndex int
	// The style for series
	Style chart.Style
	// The label for series
	Label SeriesLabel
	// The name of series
	Name string
	// Radius for Pie chart, e.g.: 40%, default is "40%"
	Radius string
	// Mark point for series
	MarkPoint SeriesMarkPoint
	// Make line for series
	MarkLine SeriesMarkLine
}
type SeriesList []Series

type PieSeriesOption struct {
	Radius string
	Label  SeriesLabel
	Names  []string
}

func NewPieSeriesList(values []float64, opts ...PieSeriesOption) []Series {
	result := make([]Series, len(values))
	var opt PieSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	for index, v := range values {
		name := ""
		if index < len(opt.Names) {
			name = opt.Names[index]
		}
		s := Series{
			Type: ChartTypePie,
			Data: []SeriesData{
				{
					Value: v,
				},
			},
			Radius: opt.Radius,
			Label:  opt.Label,
			Name:   name,
		}
		result[index] = s
	}
	return result
}

type seriesSummary struct {
	MaxIndex     int
	MaxValue     float64
	MinIndex     int
	MinValue     float64
	AverageValue float64
}

func (s *Series) Summary() seriesSummary {
	minIndex := -1
	maxIndex := -1
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64
	sum := float64(0)
	for j, item := range s.Data {
		if item.Value < minValue {
			minIndex = j
			minValue = item.Value
		}
		if item.Value > maxValue {
			maxIndex = j
			maxValue = item.Value
		}
		sum += item.Value
	}
	return seriesSummary{
		MaxIndex:     maxIndex,
		MaxValue:     maxValue,
		MinIndex:     minIndex,
		MinValue:     minValue,
		AverageValue: sum / float64(len(s.Data)),
	}
}

func (sl SeriesList) Names() []string {
	names := make([]string, len(sl))
	for index, s := range sl {
		names[index] = s.Name
	}
	return names
}

type LabelFormatter func(index int, value float64, percent float64) string

func NewPieLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	if len(layout) == 0 {
		layout = "{b}: {d}"
	}
	return NewLabelFormatter(seriesNames, layout)
}

func NewValueLabelFormater(seriesNames []string, layout string) LabelFormatter {
	if len(layout) == 0 {
		layout = "{c}"
	}
	return NewLabelFormatter(seriesNames, layout)
}

func NewLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	return func(index int, value, percent float64) string {
		// 如果无percent的则设置为<0
		percentText := ""
		if percent >= 0 {
			percentText = humanize.FtoaWithDigits(percent*100, 2) + "%"
		}
		valueText := humanize.FtoaWithDigits(value, 2)
		name := ""
		if len(seriesNames) > index {
			name = seriesNames[index]
		}
		text := strings.ReplaceAll(layout, "{c}", valueText)
		text = strings.ReplaceAll(text, "{d}", percentText)
		text = strings.ReplaceAll(text, "{b}", name)
		return text
	}
}