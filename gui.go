package main

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

type DockerApi struct {
	dockerClient *client.Client
	app          *tview.Application
	lastFocus    tview.Primitive
	grid         *tview.Grid
	gridMain     *tview.Grid
	gridDocker   *tview.Grid
	gridSwarm    *tview.Grid
	table        *tview.Table
	swarmTable   *tview.Table
	dropdown     *tview.DropDown
	list         *tview.List
	pagesTabel   *tview.Pages
	pagesMain    *tview.Pages
	text         *tview.TextView

	containerData []types.Container
	swarmData     []swarm.Service
}

const refreshInterval = 1 * time.Second

func GetApp() *tview.Application {

	app := tview.NewApplication()

	return app

}

func (d *DockerApi) RunGui() {
	d.app = GetApp()

	d.DockerGrid()
	d.GridSwarm()

	d.list = tview.NewList().
		AddItem("Docker", "docker localhost", 'a', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("dockerPage")
			d.app.SetFocus(d.table)
			d.table.SetSelectable(true, false)

		}).
		AddItem("Swarm", "Nodes info", 'b', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("swarmPage")
			d.app.SetFocus(d.swarmTable)
			d.swarmTable.SetSelectable(true, false)

		}).
		AddItem("Services", "services Info", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			d.app.Stop()

		})
	d.gridMain = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.list, 0, 0, 5, 1, 10, 10, true)

	d.pagesMain = tview.NewPages().
		AddPage("mainPage", d.gridMain, true, true).
		AddPage("dockerPage", d.gridDocker, true, false).
		AddPage("swarmPage", d.gridSwarm, true, false)

	d.MainNavigation()
	if err := d.app.SetRoot(d.pagesMain, true).Run(); err != nil {
		panic(err)
	}

}

func (d *DockerApi) DockerGrid() {

	// newPrimitive := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }
	d.table = d.ContainerTable()
	// main := newPrimitive("Main content")

	// d.list = tview.NewList().
	// 	AddItem("Docker", "docker localhost", 'a', func() {
	// 		d.lastFocus = d.app.GetFocus()
	// 		d.app.SetFocus(d.table)
	// 		d.table.SetSelectable(true, false)

	// 	}).
	// 	AddItem("Swarm", "Nodes info", 'b', func() {
	// 		d.lastFocus = d.app.GetFocus()
	// 		d.pagesMain.ShowPage("swarmPage")
	// 		d.app.SetFocus(d.swarmTable)
	// 		d.swarmTable.SetSelectable(true, false)

	// 	}).
	// 	AddItem("Services", "services Info", 'c', nil).
	// 	AddItem("List item 4", "Some explanatory text", 'd', nil).
	// 	AddItem("Quit", "Press to exit", 'q', func() {
	// 		d.app.Stop()

	// 	})

	d.dropdown = tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Restart", "Meta", "Logs", "Fourth", "Fifth"}, nil)
	textView := tview.NewTextView().SetLabel(fmt.Sprint(d.table.GetRowCount())).
		SetDynamicColors(true).
		SetRegions(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint((d.table.GetRowCount() - 1)) + " Containers")
		}
	}()
	d.text = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})

	d.pagesTabel = tview.NewPages().
		AddPage("containerTablePage", d.table, true, true).
		AddPage("containerLogInfoPage", d.text, true, false)

	d.grid = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 4, 0, 2, false).
		AddItem(d.pagesTabel, 2, 1, 4, 4, 3, 4, false)

	d.gridDocker = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		// AddItem(newPrimitive("F1"), 0, 0, 1, 1, 5, 10, false).
		AddItem(d.grid, 0, 0, 5, 1, 10, 10, false)
	// AddItem(d.list, 0, 0, 5, 1, 10, 10, true)

}
