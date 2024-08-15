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
	"errors"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
)

type pieChart struct {
	p   *Painter
	opt *PieChartOption
}

type PieChartOption struct {
	// The theme
	Theme ColorPalette
	// The font size
	Font *truetype.Font
	// The data series list
	SeriesList SeriesList
	// The padding of line chart
	Padding Box
	// The option of title
	Title TitleOption
	// The legend option
	Legend LegendOption
	// background is filled
	backgroundIsFilled bool
}

// NewPieChart returns a pie chart renderer
func NewPieChart(p *Painter, opt PieChartOption) *pieChart {
	if opt.Theme == nil {
		opt.Theme = defaultTheme
	}
	return &pieChart{
		p:   p,
		opt: &opt,
	}
}

type sector struct {
	value       float64
	percent     float64
	cx          int
	cy          int
	rx          float64
	ry          float64
	start       float64
	delta       float64
	offset      int
	quadrant    int
	lineStartX  int
	lineStartY  int
	lineBranchX int
	lineBranchY int
	lineEndX    int
	lineEndY    int
	showLabel   bool
	label       string
	series      Series
	color       Color
}

func NewSector(cx int, cy int, radius float64, labelRadius float64, value float64, currentValue float64, totalValue float64, labelLineLength int, label string, series Series, color Color) sector {
	s := sector{}
	s.value = value
	s.percent = value / totalValue
	s.cx = cx
	s.cy = cy
	s.rx = radius
	s.ry = radius
	p := (currentValue + value/2) / totalValue
	if p < 0.25 {
		s.quadrant = 1
	} else if p < 0.5 {
		s.quadrant = 4
	} else if p < 0.75 {
		s.quadrant = 3
	} else {
		s.quadrant = 2
	}
	s.start = chart.PercentToRadians(currentValue/totalValue) - math.Pi/2
	s.delta = chart.PercentToRadians(value / totalValue)
	angle := s.start + s.delta/2
	s.lineStartX = cx + int(radius*math.Cos(angle))
	s.lineStartY = cy + int(radius*math.Sin(angle))
	s.lineBranchX = cx + int(labelRadius*math.Cos(angle))
	s.lineBranchY = cy + int(labelRadius*math.Sin(angle))
	s.offset = labelLineLength
	if s.lineBranchX <= cx {
		s.offset *= -1
	}
	s.lineEndX = s.lineBranchX + s.offset
	s.lineEndY = s.lineBranchY
	s.series = series
	s.color = color
	s.showLabel = series.Label.Show
	s.label = NewPieLabelFormatter([]string{label}, series.Label.Formatter)(0, s.value, s.percent)
	return s
}

func (s *sector) calculateY(prevY int) int {
	for i := 0; i <= s.cy; i++ {
		if s.quadrant <= 2 {
			if (prevY - s.lineBranchY) > labelFontSize+5 {
				break
			}
			s.lineBranchY -= 1
		} else {
			if (s.lineBranchY - prevY) > labelFontSize+5 {
				break
			}
			s.lineBranchY += 1
		}
	}
	s.lineEndY = s.lineBranchY
	return s.lineBranchY
}

func (s *sector) calculateTextXY(textBox Box) (x int, y int) {
	textMargin := 3
	x = s.lineEndX + textMargin
	y = s.lineEndY + textBox.Height()>>1 - 1
	if s.offset < 0 {
		textWidth := textBox.Width()
		x = s.lineEndX - textWidth - textMargin
	}
	return
}

func (p *pieChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	opt := p.opt
	values := make([]float64, len(seriesList))
	total := float64(0)
	radiusValue := ""
	for index, series := range seriesList {
		if len(series.Radius) != 0 {
			radiusValue = series.Radius
		}
		value := float64(0)
		for _, item := range series.Data {
			value += item.Value
		}
		values[index] = value
		total += value
	}
	if total <= 0 {
		return BoxZero, errors.New("The sum value of pie chart should gt 0")
	}
	seriesPainter := result.seriesPainter
	cx := seriesPainter.Width() >> 1
	cy := seriesPainter.Height() >> 1

	diameter := chart.MinInt(seriesPainter.Width(), seriesPainter.Height())
	radius := getRadius(float64(diameter), radiusValue)

	labelLineWidth := 15
	if radius < 50 {
		labelLineWidth = 10
	}
	labelRadius := radius + float64(labelLineWidth)
	seriesNames := opt.Legend.Data
	if len(seriesNames) == 0 {
		seriesNames = seriesList.Names()
	}
	theme := opt.Theme

	currentValue := float64(0)

	var quadrant1, quadrant2, quadrant3, quadrant4 []sector
	for index, v := range values {
		series := seriesList[index]
		color := theme.GetSeriesColor(index)
		if index == len(values)-1 {
			if color == theme.GetSeriesColor(0) {
				color = theme.GetSeriesColor(1)
			}
		}
		s := NewSector(cx, cy, radius, labelRadius, v, currentValue, total, labelLineWidth, seriesNames[index], series, color)
		switch quadrant := s.quadrant; quadrant {
		case 1:
			quadrant1 = append([]sector{s}, quadrant1...)
		case 2:
			quadrant2 = append(quadrant2, s)
		case 3:
			quadrant3 = append([]sector{s}, quadrant3...)
		case 4:
			quadrant4 = append(quadrant4, s)
		}
		currentValue += v
	}
	sectors := append(quadrant1, quadrant4...)
	sectors = append(sectors, quadrant3...)
	sectors = append(sectors, quadrant2...)

	currentQuadrant := 0
	prevY := 0
	maxY := 0
	minY := 0
	for _, s := range sectors {
		seriesPainter.OverrideDrawingStyle(Style{
			StrokeWidth: 1,
			StrokeColor: s.color,
			FillColor:   s.color,
		})
		seriesPainter.MoveTo(s.cx, s.cy)
		seriesPainter.ArcTo(s.cx, s.cy, s.rx, s.ry, s.start, s.delta).LineTo(s.cx, s.cy).Close().FillStroke()
		if !s.showLabel {
			continue
		}
		if currentQuadrant != s.quadrant {
			if s.quadrant == 1 {
				minY = cy * 2
				maxY = 0
				prevY = cy * 2
			}
			if s.quadrant == 2 {
				if currentQuadrant != 3 {
					prevY = s.lineEndY
				} else {
					prevY = minY
				}
			}
			if s.quadrant == 3 {
				if currentQuadrant != 4 {
					prevY = s.lineEndY
				} else {
					minY = cy * 2
					maxY = 0
					prevY = 0
				}
			}
			if s.quadrant == 4 {
				if currentQuadrant != 1 {
					prevY = s.lineEndY
				} else {
					prevY = maxY
				}
			}
			currentQuadrant = s.quadrant
		}
		prevY = s.calculateY(prevY)
		if prevY > maxY {
			maxY = prevY
		}
		if prevY < minY {
			minY = prevY
		}
		seriesPainter.MoveTo(s.lineStartX, s.lineStartY)
		seriesPainter.LineTo(s.lineBranchX, s.lineBranchY)
		seriesPainter.MoveTo(s.lineBranchX, s.lineBranchY)
		seriesPainter.LineTo(s.lineEndX, s.lineEndY)
		seriesPainter.Stroke()
		textStyle := Style{
			FontColor: theme.GetTextColor(),
			FontSize:  labelFontSize,
			Font:      opt.Font,
		}
		if !s.series.Label.Color.IsZero() {
			textStyle.FontColor = s.series.Label.Color
		}
		seriesPainter.OverrideTextStyle(textStyle)
		x, y := s.calculateTextXY(seriesPainter.MeasureText(s.label))
		seriesPainter.Text(s.label, x, y)
	}
	return p.p.box, nil
}

func (p *pieChart) Render() (Box, error) {
	opt := p.opt

	renderResult, err := defaultRender(p.p, defaultRenderOption{
		Theme:      opt.Theme,
		Padding:    opt.Padding,
		SeriesList: opt.SeriesList,
		XAxis: XAxisOption{
			Show: FalseFlag(),
		},
		YAxisOptions: []YAxisOption{
			{
				Show: FalseFlag(),
			},
		},
		TitleOption:        opt.Title,
		LegendOption:       opt.Legend,
		backgroundIsFilled: opt.backgroundIsFilled,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypePie)
	return p.render(renderResult, seriesList)
}
