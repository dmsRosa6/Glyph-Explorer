package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type Entry struct {
	Name       string
	IsDir      bool
	Size       int64
	Executable bool
}

func LoadEntries(path string, showHidden bool) ([]Entry, error) {
	des, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	out := make([]Entry, 0, len(des)+1)
	if filepath.Dir(path) != path {
		out = append(out, Entry{Name: "..", IsDir: true})
	}

	for _, d := range des {
		name := d.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		e := Entry{Name: name, IsDir: d.IsDir()}
		if info, err := d.Info(); err == nil {
			e.Size = info.Size()
			e.Executable = !d.IsDir() && info.Mode()&0111 != 0
		}
		out = append(out, e)
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Name == ".." {
			return true
		}
		if out[j].Name == ".." {
			return false
		}
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})

	return out, nil
}

func OpenWithSystem(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
