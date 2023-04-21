package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (d *DockerApi) ContainerTable() *tview.Table {

	d.table = tview.NewTable().
		SetBorders(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			d.app.QueueUpdateDraw(func() {
				d.containerData = d.GetDockerContainer()
				color := tcell.ColorYellow
				d.table.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.table.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.table.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range d.containerData {
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
