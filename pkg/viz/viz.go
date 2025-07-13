package viz

import (
	"bytes"
	"fmt"
	"log"

	"github.com/richiejp/makemore/pkg/data"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type MatrixData [][]uint16

func (m MatrixData) Dims() (c, r int)   { 
	return len(m[0]), len(m)
}
func (m MatrixData) Z(c, r int) float64 { 
	return float64(m[r][c]) 
}
func (m MatrixData) X(c int) float64    { 
	return float64(c) 
}
func (m MatrixData) Y(r int) float64    { 
	return float64(r) 
}

func (m MatrixData) Label(i int) string { 
	C, _ := m.Dims()
	r := i / C
	c := i % C
	return fmt.Sprintf("%s%s\n%d", 
		data.CharIndxToString(r), 
		data.CharIndxToString(c),
		m[r][c],
	) 
}
func (m MatrixData) Len() int { 
	c, r := m.Dims()
	return c * r
}
func (m MatrixData) XY(i int) (x, y float64) { 
	C, _ := m.Dims()
	r := i / C
	c := i % C

	x = float64(c)
	y = float64(r)

	return
}

func Heatmap(N [][]uint16) []byte {
	md := MatrixData(N)

	plotter.DefaultFontSize = vg.Points(8)

	p := plot.New()
	hm := plotter.NewHeatMap(md, palette.Heat(8196, 1))
	labels, err := plotter.NewLabels(md)
	if err != nil {
		log.Fatalf("NewLabels: %v", err)
	}
	labels.Offset = vg.Point{ X: -8, Y: -8 }

	bluePalette := moreland.SmoothBlueTan().Palette(8196)
	hm.Palette = bluePalette

	// Add heatmap to plot
	p.Add(hm)
	p.Add(labels)

	// Remove axes
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{})
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{})

	buf := bytes.NewBuffer(nil)
	writerTo, err := p.WriterTo(vg.Points(float64(800)), vg.Points(float64(800)), "png")
	if err != nil {
		log.Fatalf("Could not write plot: %v", err)
	}
	writerTo.WriteTo(buf)

	return buf.Bytes()
}
