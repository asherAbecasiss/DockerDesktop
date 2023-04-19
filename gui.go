package main

import (
	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DockerApi struct {
	dockerClient *client.Client
	app          *tview.Application
	lastFocus    tview.Primitive
}

func GetApp() *tview.Application {

	app := tview.NewApplication()

	return app

}

func (d *DockerApi) RunGui() {
	d.app = GetApp()

	mainGrid := d.MainGrid()

	if err := d.app.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}

}

func (d *DockerApi) ContainerTable() *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	// lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	// cols, rows := 10, 40
	// word := 0

	data := d.GetDockerContainer()
	color := tcell.ColorYellow
	table.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
	table.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
	table.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

	color = tcell.ColorWhite
	for i, v := range data {
		for j := 0; j < 1; j++ {

			table.SetCell(i+1, j, tview.NewTableCell(v.Names[0]).SetAlign(tview.AlignCenter).SetTextColor(color))
			table.SetCell(i+1, j+1, tview.NewTableCell(v.Status).SetAlign(tview.AlignCenter).SetTextColor(color))
			table.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))

		}

	}

	// table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyEnter {
	// 		table.SetSelectable(true, false)
	// 	}

	// }).SetSelectedFunc(func(row int, column int) {
	// 	table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	// 	table.SetSelectable(true, false)

	// })

	return table
}

func (d *DockerApi) MainGrid() *tview.Grid {

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	// main := newPrimitive("Main content")
	t := d.ContainerTable()

	list := tview.NewList().
		AddItem("List item 1", "Some explanatory text", 'a', func() {
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(t)
			t.SetSelectable(true, false)

		}).
		AddItem("List item 2", "Some explanatory text", 'b', func() {
			r, _ := t.GetSelection()
			d.RestartContainerID(t.GetCell(r, 2).Text)

		}).
		AddItem("List item 3", "Some explanatory text", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			d.app.Stop()

		})

	dropdown := tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Restart", "Second", "Third", "Fourth", "Fifth"}, nil)

	grid := tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(newPrimitive("F1"), 0, 0, 1, 1, 1, 10, false).
		AddItem(dropdown, 0, 1, 1, 4, 1, 10, false).
		// AddItem(newPrimitive("Header3"), 0, 2, 1, 3, 1, 10, false).
		AddItem(list, 1, 0, 4, 1, 4, 4, true).
		AddItem(t, 1, 1, 4, 4, 5, 5, false)

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			t.SetSelectable(true, false)
			d.lastFocus = d.app.GetFocus()
			r, _ := t.GetSelection()

			d.app.SetFocus(dropdown)
			dropdown.SetLabel(t.GetCell(r, 0).Text)
			_, v := dropdown.GetCurrentOption()

			if v == "Restart" {
				d.RestartContainerID(t.GetCell(r, 2).Text)
				dropdown.SetCurrentOption(-1)
				d.app.SetFocus(list)
			}

		}
		if event.Key() == tcell.KeyESC {
			if d.lastFocus == nil {
				d.lastFocus = list

			}
			d.app.SetFocus(d.lastFocus)
		}
		return event
	})
	dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			if d.lastFocus == nil {
				d.lastFocus = list

			}
			d.app.SetFocus(d.lastFocus)
		}
		return event
	})
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlA {
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(dropdown)

		}
		if event.Key() == tcell.KeyF1 {
			d.app.SetFocus(list)

		}

		return event
	})
	return grid

}
