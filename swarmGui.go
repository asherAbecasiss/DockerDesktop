package main

import (
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
				d.swarmTable.SetCell(0, 1, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetTextColor(color))
				d.swarmTable.SetCell(0, 2, tview.NewTableCell("ID").SetAlign(tview.AlignCenter).SetTextColor(color))

				for i, v := range d.swarmData {
					for j := 0; j < 1; j++ {
						// if v.Status[:2] == "Up" {
						// 	color = tcell.ColorGreenYellow
						// } else {
						// 	color = tcell.ColorRed
						// }

						d.swarmTable.SetCell(i+1, j, tview.NewTableCell(v.Spec.Name).SetAlign(tview.AlignCenter).SetTextColor(color))

						d.swarmTable.SetCell(i+1, j+1, tview.NewTableCell(v.Meta.Version.String()).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.swarmTable.SetCell(i+1, j+2, tview.NewTableCell(v.ID).SetAlign(tview.AlignCenter).SetTextColor(color))

					}

				}
			})

		}
	}()
	return d.swarmTable
}

func (d *DockerApi) GridSwarm() {

	// newPrimitive := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	// }
	d.swarmTable = d.SwarmTable()

	// d.dropdown = tview.NewDropDown().
	// 	SetLabel("Select an option: ").
	// 	SetOptions([]string{"Restart", "Meta", "Logs", "Fourth", "Fifth"}, nil)
	// textView := tview.NewTextView().SetLabel(fmt.Sprint(d.swarmTable.GetRowCount())).
	// 	SetDynamicColors(true).
	// 	SetRegions(true)
	// go func() {
	// 	for {
	// 		time.Sleep(refreshInterval)
	// 		textView.SetLabel(fmt.Sprint((d.swarmTable.GetRowCount() - 1)) + " Containers")
	// 	}
	// }()

	d.text = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})

	// d.pagesTabel = tview.NewPages().
	// 	AddPage("containerTablePage", d.swarmTable, true, true).
	// 	AddPage("containerLogInfoPage", d.text, true, false)

	// grid := tview.NewGrid().
	// 	SetRows(1, -1).
	// 	SetColumns(-1).
	// 	SetBorders(true).
	// 	AddItem(d.dropdown, 0, 1, 1, 4, 1, 2, false).
	// 	AddItem(textView, 1, 1, 1, 4, 0, 2, false).
	// 	AddItem(d.pagesTabel, 2, 1, 4, 4, 3, 4, false)

	d.gridSwarm = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		// AddItem(newPrimitive("F1"), 0, 0, 1, 1, 5, 10, false).
		AddItem(d.swarmTable, 0, 0, 5, 1, 10, 10, false).
		AddItem(d.text, 0, 0, 5, 1, 10, 10, true)

}
