package stats

import (
	"fmt"
	"github.com/derekmu/g3n/gui"
)

// StatsTable is a panel to display performance statistics.
type StatsTable struct {
	gui.Panel
	shadersLabel  *gui.Label
	vaosLabel     *gui.Label
	buffersLabel  *gui.Label
	texturesLabel *gui.Label
	unisLabel     *gui.Label
	drawsLabel    *gui.Label
	cgosLabel     *gui.Label
}

// NewStatsTable creates and returns a pointer to a new statistics table panel
func NewStatsTable() *StatsTable {
	t := new(StatsTable)
	t.InitPanel(t, 0, 0)

	t.shadersLabel = gui.NewLabel("Shaders: -")
	t.Add(t.shadersLabel)

	t.vaosLabel = gui.NewLabel("VAOs: -")
	t.Add(t.vaosLabel)

	t.buffersLabel = gui.NewLabel("Buffers: -")
	t.Add(t.buffersLabel)

	t.texturesLabel = gui.NewLabel("Textures: -")
	t.Add(t.texturesLabel)

	t.unisLabel = gui.NewLabel("Uniforms/frame: -")
	t.Add(t.unisLabel)

	t.drawsLabel = gui.NewLabel("Draws/frame: -")
	t.Add(t.drawsLabel)

	t.cgosLabel = gui.NewLabel("CGO calls/frame: -")
	t.Add(t.cgosLabel)

	return t
}

// Update updates the table values from the specified stats table
func (t *StatsTable) Update(s *Stats) {
	y := 0

	t.shadersLabel.SetText(fmt.Sprintf("Shaders: %d", s.Glstats.Shaders))
	t.shadersLabel.FitSizeToText()
	t.shadersLabel.SetPosition(0, float32(y))
	y += t.shadersLabel.Height()

	t.vaosLabel.SetText(fmt.Sprintf("VAOs: %d", s.Glstats.Vaos))
	t.vaosLabel.FitSizeToText()
	t.vaosLabel.SetPosition(0, float32(y))
	y += t.vaosLabel.Height()

	t.buffersLabel.SetText(fmt.Sprintf("Buffers: %d", s.Glstats.Buffers))
	t.buffersLabel.FitSizeToText()
	t.buffersLabel.SetPosition(0, float32(y))
	y += t.buffersLabel.Height()

	t.texturesLabel.SetText(fmt.Sprintf("Textures: %d", s.Glstats.Textures))
	t.texturesLabel.FitSizeToText()
	t.texturesLabel.SetPosition(0, float32(y))
	y += t.texturesLabel.Height()

	t.unisLabel.SetText(fmt.Sprintf("Uniforms/frame: %d", s.Unisets))
	t.unisLabel.FitSizeToText()
	t.unisLabel.SetPosition(0, float32(y))
	y += t.unisLabel.Height()

	t.drawsLabel.SetText(fmt.Sprintf("Draws/frame: %d", s.Drawcalls))
	t.drawsLabel.FitSizeToText()
	t.drawsLabel.SetPosition(0, float32(y))
	y += t.drawsLabel.Height()

	t.cgosLabel.SetText(fmt.Sprintf("CGO calls/frame: %d", s.Cgocalls))
	t.cgosLabel.FitSizeToText()
	t.cgosLabel.SetPosition(0, float32(y))
	y += t.cgosLabel.Height()

	// assuming this is the widest label rather than checking all the label widths
	t.SetContentSize(t.cgosLabel.Width(), y)
}
