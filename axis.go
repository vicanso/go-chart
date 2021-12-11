// MIT License

// Copyright (c) 2021 Tree Xie

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
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type (
	// AxisData string
	XAxis struct {
		// 'value' 数值轴，适用于连续数据。
		// 'category' 类目轴，适用于离散的类目数据。为该类型时类目数据可自动从 series.data 或 dataset.source 中取，或者可通过 xAxis.data 设置类目数据。
		// 'time' 时间轴，适用于连续的时序数据，与数值轴相比时间轴带有时间的格式化，在刻度计算上也有所不同，例如会根据跨度的范围来决定使用月，星期，日还是小时范围的刻度。
		// 'log' 对数轴。适用于对数数据。
		Type string
		Data []string
	}
)

const axisStrokeWidth = 1

func GetXAxisAndValues(xAxis XAxis, theme string) (chart.XAxis, []float64) {
	xValues := make([]float64, len(xAxis.Data))
	ticks := make([]chart.Tick, len(xAxis.Data))
	for index, key := range xAxis.Data {
		f := float64(index)
		xValues[index] = f
		ticks[index] = chart.Tick{
			Value: f,
			Label: key,
		}
	}
	// TODO
	if theme == ThemeDark {
		return chart.XAxis{
			Ticks: ticks,
		}, xValues
	}
	return chart.XAxis{
		Ticks: ticks,
		Style: chart.Style{
			FontColor:   AxisColorLight,
			StrokeColor: AxisColorLight,
			StrokeWidth: axisStrokeWidth,
		},
	}, xValues
}

func GetYAxis(theme string) chart.YAxis {
	// TODO
	if theme == ThemeDark {
		return chart.YAxis{}
	}
	strokeColor := drawing.Color{
		R: 224,
		G: 230,
		B: 241,
		A: 255,
	}
	return chart.YAxis{

		AxisType: chart.YAxisSecondary,
		GridMajorStyle: chart.Style{
			StrokeColor: strokeColor,
			StrokeWidth: axisStrokeWidth,
		},
		GridMinorStyle: chart.Style{
			StrokeColor: strokeColor,
			StrokeWidth: axisStrokeWidth,
		},
		Style: chart.Style{
			FontColor: AxisColorLight,
			// alpha 0，隐藏
			StrokeColor: hiddenColor,
			StrokeWidth: axisStrokeWidth,
		},
	}
}
