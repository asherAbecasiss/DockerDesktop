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

	gridMain   *tview.Grid
	gridDocker *tview.Grid

	gridSwarm       *tview.Grid
	containearTable *tview.Table
	swarmTable      *tview.Table
	dropdown        *tview.DropDown
	list            *tview.List
	pagesMain       *tview.Pages
	text            *tview.TextView

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
	d.initDropdown()

	d.DockerGrid()

	d.GridSwarm()

	d.list = tview.NewList().
		AddItem("Docker", "docker localhost", 'a', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("dockerPage")

			d.app.SetFocus(d.containearTable)
			d.containearTable.SetSelectable(true, false)

		}).
		AddItem("Swarm", "Nodes info", 'b', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("swarmPage")
			d.app.SetFocus(d.swarmTable)
			d.swarmTable.SetSelectable(true, false)

		}).
		AddItem("Services", "services Info", 'c', nil).
		AddItem("RMM", "remote management", 'd', nil).
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
		AddPage("swarmPage", d.gridSwarm, true, false).
		// AddPage("containerTablePage", d.table, true, false).
		AddPage("containerLogInfoPage", tview.NewGrid().
			SetRows(1, -1).
			SetColumns(-1).
			SetBorders(true).
			AddItem(d.text, 0, 1, 4, 4, 10, 10, true), true, false)

	d.MainNavigation()
	if err := d.app.SetRoot(d.pagesMain, true).Run(); err != nil {
		panic(err)
	}

}

func (d *DockerApi) DockerGrid() {

	d.containearTable = d.ContainerTable()
	textView := tview.NewTextView().SetLabel(fmt.Sprint(d.containearTable.GetRowCount())).
		SetDynamicColors(true).
		SetRegions(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint((d.containearTable.GetRowCount() - 1)) + " Containers")
		}
	}()

	d.grid = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 4, 0, 2, false).
		AddItem(d.containearTable, 2, 1, 4, 4, 3, 4, true)

	d.gridDocker = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.grid, 0, 0, 5, 1, 10, 10, false)

}

func (d *DockerApi) initDropdown() {

	d.dropdown = tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Restart", "Meta", "Logs", "Fourth", "Fifth"}, nil)
	d.text = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})
}
