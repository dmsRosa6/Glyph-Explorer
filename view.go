package main

import (
	"fmt"

	"github.com/dmsRosa6/glyph/canvas"
	"github.com/dmsRosa6/glyph/core"
	"github.com/dmsRosa6/glyph/framework"
	"github.com/dmsRosa6/glyph/geom"
	"github.com/dmsRosa6/glyph/primitive"
	"github.com/dmsRosa6/glyph/widgets"
)

type ExplorerView struct {
	Root       *widgets.Bordered
	List       *widgets.List
	PathText   *primitive.Text
	StatusText *primitive.Text
	Hints      *primitive.Text
	Explorer   *Explorer
	Width      int
	Height     int
}

func NewExplorerView(ex *Explorer, width, height int) (*ExplorerView, error) {
	root, err := widgets.NewBox(geom.NewBounds(0, 0, width, height), widgets.BoxConfig{
		Padding: 1,
		Style:   framework.Style{Bg: core.Black, Fg: core.White},
		BorderConfig: primitive.BorderConfig{
			Thickness:   1,
			BorderStyle: primitive.Rounded,
			Style:       framework.Style{Bg: core.Transparent, Fg: core.SlateGray},
		},
	})
	if err != nil {
		return nil, err
	}

	pathText, err := primitive.NewText(&geom.Point{X: 0, Y: 0}, primitive.TextConfig{Fg: core.LightSkyBlue})
	if err != nil {
		return nil, err
	}

	statusText, err := primitive.NewText(&geom.Point{X: 0, Y: 1}, primitive.TextConfig{Fg: core.Gray})
	if err != nil {
		return nil, err
	}

	innerW := width - 4
	innerH := height - 4
	listH := innerH - 5
	if listH < 1 {
		listH = 1
	}

	list, err := widgets.NewList(geom.NewBounds(0, 3, innerW, listH), widgets.ListConfig{
		Style: framework.Style{Bg: core.Transparent, Fg: core.Transparent},
	})
	if err != nil {
		return nil, err
	}

	hints, err := primitive.NewText(&geom.Point{X: 0, Y: innerH - 1}, primitive.TextConfig{
		Value: "↑/k ↓/j  Enter/l open  ←/h up  . hidden  g/G first/last  q quit",
		Fg:    core.DimGray,
	})
	if err != nil {
		return nil, err
	}

	root.AddChild(pathText)
	root.AddChild(statusText)
	root.AddChild(list)
	root.AddChild(hints)

	v := &ExplorerView{
		Root: root, List: list, PathText: pathText,
		StatusText: statusText, Hints: hints,
		Explorer: ex, Width: innerW, Height: innerH,
	}
	v.Refresh()
	return v, nil
}

func (v *ExplorerView) Refresh() {
	v.PathText.SetValue("  " + displayPath(v.Explorer.Path, v.Width-2))
	v.StatusText.SetValue("  " + v.Explorer.Status)
	v.rebuildRows()
}

func (v *ExplorerView) rebuildRows() {
	for _, c := range v.List.Children() {
		v.List.RemoveChild(c)
	}

	_, h := v.List.Size()
	v.keepSelectedVisible(h)

	end := v.Explorer.ScrollOffset + h
	if end > len(v.Explorer.Entries) {
		end = len(v.Explorer.Entries)
	}

	for i := v.Explorer.ScrollOffset; i < end; i++ {
		row, err := buildRow(v.Explorer.Entries[i], v.Width, i == v.Explorer.Selected)
		if err == nil {
			v.List.AddChild(row)
		}
	}
}

func (v *ExplorerView) keepSelectedVisible(viewH int) {
	if viewH <= 0 {
		return
	}

	s := v.Explorer.Selected
	if s < v.Explorer.ScrollOffset {
		v.Explorer.ScrollOffset = s
	} else if s >= v.Explorer.ScrollOffset+viewH {
		v.Explorer.ScrollOffset = s - viewH + 1
	}

	maxOffset := len(v.Explorer.Entries) - viewH
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.Explorer.ScrollOffset > maxOffset {
		v.Explorer.ScrollOffset = maxOffset
	}
	if v.Explorer.ScrollOffset < 0 {
		v.Explorer.ScrollOffset = 0
	}
}

func buildRow(e Entry, width int, selected bool) (*canvas.Container, error) {
	fg := core.White
	if e.Name == ".." {
		fg = core.SlateGray
	} else if e.IsDir {
		fg = core.LightGray
	} else if e.Executable {
		fg = core.LimeGreen
	}
	if selected {
		fg = core.White
	}

	bg := core.Transparent
	if selected {
		bg = core.DarkSlateGray
	}

	row, err := canvas.NewContainer(geom.NewBounds(0, 0, width, 1), canvas.ContainerConfig{
		Style: framework.Style{Bg: bg, Fg: core.Transparent},
	})
	if err != nil {
		return nil, err
	}

	fill, err := primitive.NewRect(geom.NewBounds(0, 0, width, 1), primitive.RectConfig{
		Style: framework.Style{Bg: bg, Fg: core.Transparent},
	})
	if err != nil {
		return nil, err
	}
	row.AddChild(fill)

	name := e.Name
	icon := "  "
	if e.Name == ".." {
		icon = "↩ "
	} else if e.IsDir {
		icon = "▸ "
	} else if e.Executable {
		icon = "◆ "
	} else {
		icon = "· "
	}
	if e.IsDir && name != ".." {
		name += "/"
	}

	right := ""
	if !e.IsDir {
		right = humanSize(e.Size)
	}

	avail := width - 4 - len(right)
	if avail < 1 {
		avail = 1
	}
	name = truncateRunes(name, avail)

	label, err := primitive.NewText(&geom.Point{X: 1, Y: 0},
		primitive.TextConfig{Value: icon + name, Fg: fg})
	if err != nil {
		return nil, err
	}
	row.AddChild(label)

	if right != "" {
		rx := width - len(right) - 1
		if rx > len([]rune(name))+4 {
			if sizeLabel, err := primitive.NewText(&geom.Point{X: rx, Y: 0},
				primitive.TextConfig{Value: right, Fg: core.Gray}); err == nil {
				row.AddChild(sizeLabel)
			}
		}
	}

	return row, nil
}

func humanSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	v := float64(n)
	for _, unit := range units {
		v /= 1024
		if v < 1024 {
			return fmt.Sprintf("%.1f %s", v, unit)
		}
	}
	return fmt.Sprintf("%.1f PiB", v/1024)
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max == 1 {
		return string(r[:1])
	}
	return string(r[:max-1]) + "…"
}

func displayPath(path string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	runes := []rune(path)
	if len(runes) <= maxWidth {
		return path
	}
	if maxWidth == 1 {
		return string(runes[len(runes)-1:])
	}

	return "…" + string(runes[len(runes)-(maxWidth-1):])
}
