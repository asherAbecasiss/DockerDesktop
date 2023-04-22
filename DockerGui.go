package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (d *DockerApi) DockerGrid() {

	d.containearTable = d.ContainerTable()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	textView2 := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	textView3 := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint((d.containearTable.GetRowCount() - 1)) + " Containers")

		}
	}()

	fmt.Fprintf(textView2, " [yellow]Filters[white] \n [green]F1[white]: All \n [green]F2[white]: Status up")
	fmt.Fprintf(textView3, " %d", *d.filters)
	d.grid = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 1, 0, 2, false).
		AddItem(textView2, 1, 2, 1, 1, 0, 2, false).
		AddItem(textView3, 1, 3, 1, 2, 0, 2, false).
		AddItem(d.containearTable, 2, 1, 4, 4, 3, 4, true)

	d.gridDocker = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.grid, 0, 0, 5, 1, 10, 10, false)

}
func (d *DockerApi) DockerGridImage() {

	d.GetImageList()
	textView := tview.NewTextView().SetLabel(fmt.Sprint(d.containearTableImage.GetRowCount())).
		SetDynamicColors(true).
		SetRegions(true)

	go func() {
		for {
			time.Sleep(refreshInterval)
			textView.SetLabel(fmt.Sprint((d.containearTableImage.GetRowCount() - 1)) + " Images")
		}
	}()

	grid := tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdownImageList, 0, 1, 1, 4, 1, 2, false).
		AddItem(textView, 1, 1, 1, 4, 0, 2, false).
		AddItem(d.containearTableImage, 2, 1, 4, 4, 3, 4, true)

	d.gridDockerImage = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(grid, 0, 0, 5, 1, 10, 10, false)

}
func (d *DockerApi) ContainerTable() *tview.Table {

	d.containearTable = tview.NewTable().
		SetBorders(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			d.app.QueueUpdateDraw(func() {
				d.containerData = d.GetDockerContainer()

				color := tcell.ColorYellow
				d.containearTable.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.containearTable.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.containearTable.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range d.containerData {

					for j := 0; j < 1; j++ {

						if *d.filters == 1 {
							if v.Status[:2] == "Up" {
								color = tcell.ColorGreenYellow
							} else {
								color = tcell.ColorRed
							}

							d.containearTable.SetCell(i+1, j, tview.NewTableCell(v.Names[0]).SetAlign(tview.AlignCenter).SetTextColor(color))

							d.containearTable.SetCell(i+1, j+1, tview.NewTableCell(v.Status).SetAlign(tview.AlignCenter).SetTextColor(color))
							d.containearTable.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))

						}
						if *d.filters == 2 {
							if v.Status[:2] == "Up" {
								color = tcell.ColorGreenYellow

								d.containearTable.SetCell(i+1, j, tview.NewTableCell(v.Names[0]).SetAlign(tview.AlignCenter).SetTextColor(color))

								d.containearTable.SetCell(i+1, j+1, tview.NewTableCell(v.Status).SetAlign(tview.AlignCenter).SetTextColor(color))
								d.containearTable.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))
							}

						}

					}

				}
			})

		}
	}()

	return d.containearTable
}

func (d *DockerApi) GetImageList() {
	d.containearTableImage = tview.NewTable().
		SetBorders(true)
	go func() {
		for {
			time.Sleep(refreshInterval)
			d.app.QueueUpdateDraw(func() {
				d.imageList = d.GetImageListDocker()
				color := tcell.ColorYellow
				d.containearTableImage.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.containearTableImage.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.containearTableImage.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range d.imageList {
					for j := 0; j < 1; j++ {

						d.containearTableImage.SetCell(i+1, j, tview.NewTableCell(v.RepoDigests).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.containearTableImage.SetCell(i+1, j+1, tview.NewTableCell(v.RepoTags).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.containearTableImage.SetCell(i+1, j+2, tview.NewTableCell(v.Id).SetAlign(tview.AlignCenter).SetTextColor(color))
					}

				}
			})

		}
	}()

}
