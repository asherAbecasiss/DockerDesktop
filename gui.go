package main

import (
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DockerApi struct {
	dockerClient *client.Client
	app          *tview.Application
	lastFocus    tview.Primitive
	grid         *tview.Grid
	table        *tview.Table
	dropdown     *tview.DropDown
	list         *tview.List
	pages        *tview.Pages
	box          *tview.Box
	modal        *tview.Modal
}

const refreshInterval = 500 * time.Millisecond

func GetApp() *tview.Application {

	app := tview.NewApplication()

	return app

}

func (d *DockerApi) RunGui() {
	d.app = GetApp()

	d.grid = d.MainGrid()
	d.MainNavigation()
	d.modal = tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {

				d.pages.SwitchToPage("main")
				d.app.SetFocus(d.dropdown)
			}
		})
	d.box = tview.NewBox().
		SetBorder(true).
		SetTitle("Centered Box")
	d.pages = tview.NewPages().
		AddPage("main", d.grid, true, true).
		AddPage("modal", d.modal, true, false)

	if err := d.app.SetRoot(d.pages, true).Run(); err != nil {
		panic(err)
	}

}

func (d *DockerApi) ContainerTable() *tview.Table {

	d.table = tview.NewTable().
		SetBorders(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			d.app.QueueUpdateDraw(func() {
				data := d.GetDockerContainer()
				color := tcell.ColorYellow
				d.table.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.table.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.table.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range data {
					for j := 0; j < 1; j++ {
						if v.Status[:2] == "Up" {
							color = tcell.ColorGreenYellow
						} else {
							color = tcell.ColorRed
						}

						d.table.SetCell(i+1, j, tview.NewTableCell(v.Names[0]).SetAlign(tview.AlignCenter).SetTextColor(color))

						d.table.SetCell(i+1, j+1, tview.NewTableCell(v.Status).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.table.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))

					}

				}
			})

		}
	}()

	return d.table
}

func (d *DockerApi) MainGrid() *tview.Grid {

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	d.table = d.ContainerTable()
	// main := newPrimitive("Main content")

	d.list = tview.NewList().
		AddItem("Docker", "docker localhost", 'a', func() {
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(d.table)
			d.table.SetSelectable(true, false)

		}).
		AddItem("Swarm", "Nodes info", 'b', func() {
			r, _ := d.table.GetSelection()
			d.RestartContainerID(d.table.GetCell(r, 2).Text)

		}).
		AddItem("Services", "services Info", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			d.app.Stop()

		})

	d.dropdown = tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Restart", "Meta", "Third", "Fourth", "Fifth"}, nil)
	textView := tview.NewTextView().SetLabel(fmt.Sprint(d.table.GetRowCount())).
		SetDynamicColors(true).
		SetRegions(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint(d.table.GetRowCount()) + " Containers")
		}
	}()

	d.grid = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(newPrimitive("F1"), 0, 0, 1, 1, 1, 1, false).
		AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 4, 0, 2, false).
		AddItem(d.list, 1, 0, 5, 1, 10, 10, true).
		AddItem(d.table, 2, 1, 4, 4, 3, 4, false)

	return d.grid

}

func (d *DockerApi) MainNavigation() {

	d.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			d.dropdown.SetCurrentOption(-1)
			d.table.SetSelectable(true, false)
			d.lastFocus = d.app.GetFocus()
			r, _ := d.table.GetSelection()

			d.app.SetFocus(d.dropdown)
			d.dropdown.SetLabel(d.table.GetCell(r, 0).Text)

		}
		if event.Key() == tcell.KeyESC {
			if d.lastFocus == nil {
				d.lastFocus = d.list

			}

			d.app.SetFocus(d.lastFocus)
		}
		return event
	})
	d.dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			if d.lastFocus == nil {
				d.lastFocus = d.table

			}
			d.dropdown.SetLabel("Select an option: ")
			d.dropdown.SetCurrentOption(-1)

			d.app.SetFocus(d.lastFocus)
		}
		if event.Key() == tcell.KeyEnter {

			i, _ := d.dropdown.GetCurrentOption()
			r, _ := d.table.GetSelection()

			if i == 0 {

				d.RestartContainerID(d.table.GetCell(r, 2).Text)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.app.SetFocus(d.lastFocus)

				return nil
			}
			if i == 1 {

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.pages.ShowPage("modal")
				d.modal.SetText(d.table.GetCell(r, 2).Text)
				d.app.SetFocus(d.modal)
				return nil

			}

		}
		return event
	})

	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlA {
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(d.dropdown)

		}
		if event.Key() == tcell.KeyF1 {
			d.app.SetFocus(d.list)

		}
		if event.Key() == tcell.KeyEscape {
			d.pages.SwitchToPage("main")

		}

		return event
	})

}
