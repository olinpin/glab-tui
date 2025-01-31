package main

import (
	"github.com/rivo/tview"
)

type Issue struct {
    id int
    text string
}

func main() {
	// box := tview.
	// 	NewBox().
	// 	SetBorder(true).
	// 	SetTitle("Hello, world!")
	app := tview.NewApplication()
	helpList := grid(app)
	if err := app.SetRoot(helpList, true).SetFocus(helpList).Run(); err != nil {
		panic(err)
	}
	// if err := app.SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }
}

func grid(app *tview.Application) *tview.Flex {

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewBox().
			SetTitle(text).
			SetBorder(true)

	}

	menu := newPrimitive("Menu")
	main := newPrimitive("Main Content")

	grid := tview.NewFlex().
		// SetRows(3, 0, 3).
		// SetColumns(30, 0, 30).
		// SetBorders(true).
		AddItem(menu, 45, 1, true).
		AddItem(main, 0, 3, false)

	return grid
}

func help(app *tview.Application) *tview.List {
	list := tview.NewList().
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		}).
		AddItem("List item 2", "Explain", 'j', func() {
		}).
		AddItem("List item 3", "Explain", 'k', func() {
		})
	return list
}
