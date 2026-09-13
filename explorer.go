package main

import (
	"fmt"
	"path/filepath"
)

type Explorer struct {
	Path         string
	Entries      []Entry
	Selected     int
	ScrollOffset int
	ShowHidden   bool
	Status       string
}

func NewExplorer(path string) *Explorer {
	e := &Explorer{}
	e.SetPath(path)
	return e
}

func (e *Explorer) SetPath(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		e.Status = err.Error()
		return
	}

	entries, err := LoadEntries(abs, e.ShowHidden)
	if err != nil {
		e.Status = fmt.Sprintf("cannot open %s: %v", filepath.Base(abs), err)
		return
	}

	e.Path = abs
	e.Entries = entries
	e.Selected = 0
	e.ScrollOffset = 0
	e.Status = fmt.Sprintf("%d items", len(entries))
}

func (e *Explorer) Reload() {
	selectedName := ""
	if e.Selected >= 0 && e.Selected < len(e.Entries) {
		selectedName = e.Entries[e.Selected].Name
	}

	entries, err := LoadEntries(e.Path, e.ShowHidden)
	if err != nil {
		e.Status = err.Error()
		return
	}

	e.Entries = entries
	e.Selected = 0
	for i, entry := range entries {
		if entry.Name == selectedName {
			e.Selected = i
			break
		}
	}
	e.Status = fmt.Sprintf("%d items", len(entries))
}

func (e *Explorer) MoveSelection(delta int) {
	if len(e.Entries) == 0 {
		return
	}
	e.Selected += delta
	if e.Selected < 0 {
		e.Selected = 0
	}
	if e.Selected >= len(e.Entries) {
		e.Selected = len(e.Entries) - 1
	}
}

func (e *Explorer) SelectFirst() {
	if len(e.Entries) > 0 {
		e.Selected = 0
	}
}

func (e *Explorer) SelectLast() {
	if len(e.Entries) > 0 {
		e.Selected = len(e.Entries) - 1
	}
}

func (e *Explorer) GoUp() {
	parent := filepath.Dir(e.Path)
	if parent != e.Path {
		e.SetPath(parent)
	}
}

func (e *Explorer) Activate() {
	if len(e.Entries) == 0 {
		return
	}

	entry := e.Entries[e.Selected]
	if entry.Name == ".." {
		e.GoUp()
		return
	}

	full := filepath.Join(e.Path, entry.Name)
	if entry.IsDir {
		e.SetPath(full)
		return
	}

	if err := OpenWithSystem(full); err != nil {
		e.Status = fmt.Sprintf("could not open %s: %v", entry.Name, err)
		return
	}
	e.Status = fmt.Sprintf("opened %s", entry.Name)
}

func (e *Explorer) ToggleHidden() {
	e.ShowHidden = !e.ShowHidden
	e.Reload()
	if e.ShowHidden {
		e.Status = "hidden files enabled"
	} else {
		e.Status = "hidden files disabled"
	}
}
