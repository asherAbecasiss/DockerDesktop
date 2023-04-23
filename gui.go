package main

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

type DockerApi struct {
	dockerClient         *client.Client
	app                  *tview.Application
	lastFocus            tview.Primitive
	grid                 *tview.Grid
	tableProcesses       *tview.Table
	gridMain             *tview.Grid
	gridDocker           *tview.Grid
	gridDockerImage      *tview.Grid
	gridDashboard        *tview.Grid
	gridSwarm            *tview.Grid
	containearTable      *tview.Table
	containearTableImage *tview.Table
	swarmTable           *tview.Table
	dropdown             *tview.DropDown
	dropdownImageList    *tview.DropDown
	list                 *tview.List
	pagesMain            *tview.Pages
	text                 *tview.TextView

	containerData            []types.Container
	swarmData                []swarm.Service
	imageList                []ImageTag
	filters                  *int
	refreshIntervalProcesses *int
	quit                     chan bool
	flagForPs                int
}

const refreshInterval = 1 * time.Second

func GetApp() *tview.Application {

	app := tview.NewApplication()

	return app

}

func (d *DockerApi) InitOpt() {
	d.filters = new(int)
	*d.filters = 1

	d.refreshIntervalProcesses = new(int)
	*d.refreshIntervalProcesses = 1

	d.quit = make(chan bool)

	d.flagForPs = 0
}

func (d *DockerApi) RunGui() {
	d.app = GetApp()
	d.InitOpt()
	d.initDropdown()
	d.DashboardGrid()

	d.DockerGrid()
	d.DockerGridImage()
	d.GridSwarm()

	d.list = tview.NewList().
		AddItem("Dashboard", "System Information", 'f', func() {
			d.lastFocus = d.app.GetFocus()

			d.pagesMain.ShowPage("dashboardPage")
			d.app.SetFocus(d.tableProcesses)
			// d.flagForPs = 1

		}).
		AddItem("Docker", "docker localhost", 'a', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("dockerPage")
			d.pagesMain.HidePage("Image")
			d.app.SetFocus(d.containearTable)
			d.containearTable.SetSelectable(true, false)

		}).
		AddItem("Docker Images", "Image List", 'b', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("imageList")
			d.pagesMain.HidePage("Image")
			d.app.SetFocus(d.containearTableImage)
			d.containearTableImage.SetSelectable(true, false)

		}).
		AddItem("Swarm", "Nodes info", 'c', func() {
			d.lastFocus = d.app.GetFocus()
			d.pagesMain.ShowPage("swarmPage")
			d.pagesMain.HidePage("Image")
			d.app.SetFocus(d.swarmTable)
			d.swarmTable.SetSelectable(true, false)

		}).
		AddItem("Services", "services Info", 'd', nil).
		AddItem("RMM", "remote management", 'e', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			d.app.Stop()

		})
	d.gridMain = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.list, 0, 0, 5, 1, 10, 10, true)

	image := tview.NewImage()
	b, _ := base64.StdEncoding.DecodeString(beach)
	photo, _ := jpeg.Decode(bytes.NewReader(b))

	image.SetImage(photo)

	//---------------------------------------------------
	d.pagesMain = tview.NewPages().
		AddPage("mainPage", d.gridMain, true, true).
		AddPage("dockerPage", d.gridDocker, true, false).
		AddPage("imageList", d.gridDockerImage, true, false).
		AddPage("swarmPage", d.gridSwarm, true, false).
		AddPage("dashboardPage", d.gridDashboard, true, false).

		// AddPage("containerTablePage", d.table, true, false).
		AddPage("containerLogInfoPage", tview.NewGrid().
			SetRows(1, -1).
			SetColumns(-1).
			SetBorders(true).
			AddItem(d.text, 0, 1, 4, 4, 10, 10, true), true, false)
	// 	AddPage("Image", tview.NewGrid().
	// 		SetRows(1, -1).
	// 		SetColumns(-1).
	// 		SetBorders(true).
	// 		AddItem(image, 0, 1, 4, 4, 10, 10, true), true, false)
	// d.pagesMain.ShowPage("Image")
	// d.pagesMain.SetChangedFunc(func() {
	// 	j, _ := d.pagesMain.GetFrontPage()
	// 	if j == "mainPage" {
	// 		d.pagesMain.ShowPage("Image")
	// 	}
	// })

	d.MainNavigation()
	if err := d.app.SetRoot(d.pagesMain, true).Run(); err != nil {
		panic(err)
	}

}

func (d *DockerApi) initDropdown() {

	d.dropdown = tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Restart", "Meta", "Logs", "Stop", "Start"}, nil)

	d.dropdownImageList = tview.NewDropDown().
		SetLabel("Select an option: ").
		SetOptions([]string{"Delete"}, nil)
	d.text = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})
}
