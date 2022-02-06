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
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type BarStyle struct {
	ClassName       string
	StrokeDashArray []float64
	FillColor       drawing.Color
}

func (bs *BarStyle) Style() chart.Style {
	return chart.Style{
		ClassName:       bs.ClassName,
		StrokeDashArray: bs.StrokeDashArray,
		StrokeColor:     bs.FillColor,
		StrokeWidth:     1,
		FillColor:       bs.FillColor,
	}
}

// Bar renders bar for chart
func (d *Draw) Bar(b chart.Box, style BarStyle) {
	s := style.Style()

	r := d.Render
	s.GetFillAndStrokeOptions().WriteToRenderer(r)
	d.moveTo(b.Left, b.Top)
	d.lineTo(b.Right, b.Top)
	d.lineTo(b.Right, b.Bottom)
	d.lineTo(b.Left, b.Bottom)
	d.lineTo(b.Left, b.Top)
	d.Render.FillStroke()
}