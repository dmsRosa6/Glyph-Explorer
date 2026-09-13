package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmsRosa6/glyph/app"
	"github.com/dmsRosa6/glyph/core"
	"github.com/dmsRosa6/glyph/framework"
	"github.com/dmsRosa6/glyph/render"
	"github.com/dmsRosa6/glyph/term"
)

func main() {
	start := "."
	if len(os.Args) > 1 {
		start = os.Args[1]
	}

	abs, err := filepath.Abs(start)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	info, err := os.Stat(abs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		abs = filepath.Dir(abs)
	}

	size, err := term.TermSize()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not determine terminal size:", err)
		os.Exit(1)
	}

	width, height := size.Cols-1, size.Rows-1
	if width < 30 {
		width = 30
	}
	if height < 12 {
		height = 12
	}

	a, err := app.NewApp(app.AppConfig{
		Width:          width,
		Height:         height,
		Bg:             core.Black,
		Fg:             core.White,
		RenderMode:     render.OnDemandMode(),
		DisableFileLog: true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ex := NewExplorer(abs)
	view, err := NewExplorerView(ex, width, height)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	a.Canvas.AddShape(view.Root)
	bindKeys(a, ex, view)
	a.Run()
}

func bindKeys(a *app.App, ex *Explorer, view *ExplorerView) {
	a.BindKey(framework.KeyUp, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.MoveSelection(-1)
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyDown, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.MoveSelection(1)
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyLeft, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.GoUp()
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyRight, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.Activate()
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyBackspace, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.GoUp()
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyEnter, func(_ framework.AppContext, _ framework.Event) (bool, error) {
		ex.Activate()
		view.Refresh()
		return true, nil
	})
	a.BindKey(framework.KeyRune, func(ctx framework.AppContext, ev framework.Event) (bool, error) {
		switch ev.Rune {
		case 'q':
			ctx.SignalApp(core.SIGTERM)
		case 'j':
			ex.MoveSelection(1)
		case 'k':
			ex.MoveSelection(-1)
		case 'l':
			ex.Activate()
		case 'h':
			ex.GoUp()
		case 'g':
			ex.SelectFirst()
		case 'G':
			ex.SelectLast()
		case '.':
			ex.ToggleHidden()
		}
		view.Refresh()
		return true, nil
	})
}
