package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//Swarm

func (d *DockerApi) SwarmTable() *tview.Table {

	d.swarmTable = tview.NewTable().
		SetBorders(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			d.app.QueueUpdateDraw(func() {
				d.swarmData = d.GetDockerServices()

				color := tcell.ColorYellow
				d.swarmTable.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.swarmTable.SetCell(0, 1, tview.NewTableCell("Image").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.swarmTable.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range d.swarmData {
					for j := 0; j < 1; j++ {
						// if v.Status[:2] == "Up" {
						// 	color = tcell.ColorGreenYellow
						// } else {
						// 	color = tcell.ColorRed
						// }

						d.swarmTable.SetCell(i+1, j, tview.NewTableCell(v.Spec.Name).SetAlign(tview.AlignCenter).SetTextColor(color))

						d.swarmTable.SetCell(i+1, j+1, tview.NewTableCell(v.Spec.TaskTemplate.ContainerSpec.Image).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.swarmTable.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))

					}

				}
			})

		}
	}()
	return d.swarmTable
}

func (d *DockerApi) GridSwarm() {

	d.swarmTable = d.SwarmTable()

	textView := tview.NewTextView().SetLabel(fmt.Sprint(d.swarmTable.GetRowCount())).
		SetDynamicColors(true).
		SetRegions(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint((d.swarmTable.GetRowCount() - 1)) + " Services")
		}
	}()

	d.grid = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 4, 0, 2, false).
		AddItem(d.swarmTable, 2, 1, 4, 4, 3, 4, true)

	d.gridSwarm = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.grid, 0, 0, 5, 1, 10, 10, false)

}
