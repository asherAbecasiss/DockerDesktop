package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (d *DockerApi) GetPsGoFunc() {
	color := tcell.ColorWhite
	for {
		select {
		case <-d.quit:
			return
		default:
			d.app.QueueUpdateDraw(func() {
				res := GetTotalProcesses()
				d.tableProcesses.SetCell(0, 0, tview.NewTableCell("Name").SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow))
				d.tableProcesses.SetCell(0, 1, tview.NewTableCell("Pid").SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow))
				d.tableProcesses.SetCell(0, 2, tview.NewTableCell("CpuPercent").SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreenYellow))
				var sum float64
				for i, v := range res {

					for j := 0; j < 1; j++ {

						switch {
						case v.CpuPercent > 10:
							color = tcell.ColorRed
						case v.CpuPercent > 1 && v.CpuPercent < 10:
							color = tcell.ColorAqua
						default:
							color = tcell.ColorWhite
						}
						sum += v.CpuPercent
						d.tableProcesses.SetCell(i+1, j, tview.NewTableCell(v.Name).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.tableProcesses.SetCell(i+1, j+1, tview.NewTableCell(fmt.Sprint(v.Id)).SetAlign(tview.AlignCenter).SetTextColor(color))
						d.tableProcesses.SetCell(i+1, j+2, tview.NewTableCell(fmt.Sprintf("%.2f", float32(v.CpuPercent))).SetAlign(tview.AlignCenter).SetTextColor(color))

					}

				}
				// d.memoryText.SetLabel(fmt.Sprint((d.tableProcesses.GetRowCount() - 1)) + " Processes\n")

				memModel := ReadMemoryStats()
				temprature := GetTemperatureStat()
				disk := GetDiskServices("/")

				d.memoryText.Clear()

				memModel.MemAvailable = memModel.MemAvailable / 1000000
				memModel.MemFree = memModel.MemFree / 1000000
				memModel.MemTotal = memModel.MemTotal / 1000000
				memModel.MemPercent = int((100 - (memModel.MemAvailable/memModel.MemTotal)*100))
				colMem := "lime"
				colDisk := "lime"
				if memModel.MemPercent > 80 {
					colMem = "red"
				}

				if (disk.Total - disk.Used) < 30 {
					colDisk = "red"
				}
				fmt.Fprintf(d.memoryText, " [yellow]Processes[white] %d \n", (d.tableProcesses.GetRowCount() - 1))
				fmt.Fprintf(d.memoryText, " [yellow]Total CPU Percent[white] %f %%  1 core=100\n", float32(sum))

				fmt.Fprintf(d.memoryText, " [navy]Free Memory[white] ["+colMem+"] %f gb\n", memModel.MemAvailable)
				fmt.Fprintf(d.memoryText, " [navy]Total Memory[white] %f gb\n", memModel.MemTotal)
				fmt.Fprintf(d.memoryText, " [navy]MemPercent Memory[white] ["+colMem+"] %d%% \n", memModel.MemPercent)
				fmt.Fprintf(d.memoryText, " [purple]Total Disk[white] %f \n", float32(disk.Total)/1000000000)
				fmt.Fprintf(d.memoryText, " [purple]Free Disk ["+colDisk+"] %f \n", float32(disk.Free)/1000000000)
				fmt.Fprintf(d.memoryText, " [purple]Used Disk[white] ["+colDisk+"] %f \n", float32(disk.Used)/1000000000)

				for _, v := range temprature {
					fmt.Fprintf(d.sensorsTemperaturesText, " [yellow]Temperature[white] %s  %f\n", v.SensorKey, v.Temperature)
				}
			})

			time.Sleep(refreshInterval)

		}
	}
}

func (d *DockerApi) TotalProcessesGui() {

	d.tableProcesses = tview.NewTable()
	d.memoryText = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	d.sensorsTemperaturesText = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	// d.tableProcesses.Clear()
	go d.GetPsGoFunc()

}

func (d *DockerApi) DashboardGrid() {

	d.TotalProcessesGui()
	textView2 := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	fmt.Fprintf(textView2, " [yellow]Processes Action[white] \n [green]Back[white]: ESC \n [green]F1[white]: Start Interval \n [green]F2[white]: Stop Interval")

	textIP := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	localIp := GetLocalIP()

	fmt.Fprintf(textIP, " [yellow]Local Ip[white] %s \n", localIp.LocalIp)

	memText2 := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	// go func() {
	// 	for {

	// 		idle0, total0 := getCPUSample()
	// 		time.Sleep(3 * time.Second)
	// 		memText2.Clear()
	// 		idle1, total1 := getCPUSample()
	// 		idleTicks := float64(idle1 - idle0)
	// 		totalTicks := float64(total1 - total0)
	// 		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	// 		var res types.Cpu

	// 		res.CpuUsage = cpuUsage
	// 		res.Busy = totalTicks - idleTicks
	// 		res.Total = totalTicks
	// 		fmt.Fprintf(memText2, " [navy]Total Cpu[white] %f \n", res.Total)
	// 		fmt.Fprintf(memText2, " [navy]Cpu Usage[white] %f \n", res.CpuUsage)
	// 		fmt.Fprintf(memText2, " [navy]Cpu Busy[white] %f \n", res.Busy)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	grid := tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(d.dropdownPS, 0, 0, 1, 7, 0, 2, false).
		AddItem(textIP, 0, 10, 1, 7, 0, 2, false).
		AddItem(d.sensorsTemperaturesText, 1, 10, 2, 10, 0, 2, false).
		AddItem(d.memoryText, 1, 0, 1, 10, 0, 2, false).
		AddItem(memText2, 2, 0, 1, 10, 0, 2, false)
	// AddItem(newPrimitive("Adding Costome Script!"), 2, 0, 4, 20, 0, 2, true)

	flex := tview.NewFlex().
		AddItem(d.tableProcesses, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			// AddItem(d.memoryText, 3, 5, false).
			AddItem(grid, 0, 5, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(textView2, 20, 1, false)

	d.gridDashboard = tview.NewGrid().
		SetRows(1, -1).
		SetColumns(-1).
		SetBorders(true).
		AddItem(flex, 0, 0, 5, 1, 10, 10, true)

}
