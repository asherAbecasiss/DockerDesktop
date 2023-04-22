package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell/v2"
)

func (d *DockerApi) MainNavigation() {

	d.SwarmNavigation()
	d.ContainearTableNavigation()
	d.DropDownNavigation()
	d.ImageListTableNavigation()
	d.DropDownImageListNavigation()

	d.text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			d.pagesMain.HidePage("containerLogInfoPage")
			d.pagesMain.ShowPage("dockerPage")
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(d.containearTable)
			return nil

		}
		return event
	})

	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlA {
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(d.dropdown)

		}
		// if event.Key() == tcell.KeyF1 {
		// 	d.app.SetFocus(d.list)

		// }
		// if event.Key() == tcell.KeyEscape {
		// 	// d.pagesTabel.ShowPage("main1")
		// 	d.lastFocus = d.app.GetFocus()
		// 	d.app.SetFocus(d.lastFocus)

		// }

		return event
	})

}

func (d *DockerApi) SwarmNavigation() {
	d.swarmTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			d.dropdown.SetCurrentOption(-1)
			d.swarmTable.SetSelectable(true, false)
			d.lastFocus = d.app.GetFocus()
			r, _ := d.swarmTable.GetSelection()

			d.app.SetFocus(d.dropdown)
			d.dropdown.SetLabel(d.swarmTable.GetCell(r, 0).Text)
			return nil

		}
		if event.Key() == tcell.KeyESC {
			d.lastFocus = d.app.GetFocus()

			d.pagesMain.SwitchToPage("mainPage")
			d.app.SetFocus(d.list)
			return nil
		}
		return event
	})
}

func (d *DockerApi) ContainearTableNavigation() {
	d.containearTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			d.dropdown.SetCurrentOption(-1)
			d.containearTable.SetSelectable(true, false)
			d.lastFocus = d.app.GetFocus()
			r, _ := d.containearTable.GetSelection()

			d.app.SetFocus(d.dropdown)
			d.dropdown.SetLabel(d.containearTable.GetCell(r, 0).Text)
			return nil

		}
		if event.Key() == tcell.KeyESC {
			d.lastFocus = d.app.GetFocus()

			d.pagesMain.SwitchToPage("mainPage")
			d.app.SetFocus(d.list)
			return nil
		}
		if event.Key() == tcell.KeyF1 {
			d.containearTable.Clear()
			*d.filters = 1

			return nil
		}
		if event.Key() == tcell.KeyF2 {
			d.containearTable.Clear()
			*d.filters = 2

			return nil
		}

		return event
	})
}

func (d *DockerApi) ImageListTableNavigation() {
	d.containearTableImage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			d.dropdownImageList.SetCurrentOption(-1)
			d.containearTableImage.SetSelectable(true, false)
			d.lastFocus = d.app.GetFocus()
			r, _ := d.containearTableImage.GetSelection()

			d.app.SetFocus(d.dropdownImageList)
			d.dropdownImageList.SetLabel(d.containearTableImage.GetCell(r, 0).Text)
			return nil

		}
		if event.Key() == tcell.KeyESC {
			d.lastFocus = d.app.GetFocus()

			d.pagesMain.SwitchToPage("mainPage")
			d.app.SetFocus(d.list)
			return nil
		}
		return event
	})
}

func (d *DockerApi) DropDownNavigation() {
	d.dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			// if d.lastFocus == nil {
			// 	d.lastFocus = d.list

			// }
			d.dropdown.SetLabel("Select an option: ")
			d.dropdown.SetCurrentOption(-1)
			d.pagesMain.HidePage("containerLogInfoPage")
			// d.lastFocus = d.containearTable
			d.app.SetFocus(d.lastFocus)

			return nil

		}
		if event.Key() == tcell.KeyEnter {

			i, _ := d.dropdown.GetCurrentOption()
			r, _ := d.containearTable.GetSelection()
			// d.lastFocus = d.app.GetFocus()
			if i == 0 {

				d.RestartContainerID(d.containearTable.GetCell(r, 2).Text)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.app.SetFocus(d.lastFocus)

				return nil
			}
			if i == 1 {
				d.text.Clear()

				d.dropdown.SetLabel("Select an option2: ")
				d.dropdown.SetCurrentOption(-1)
				d.pagesMain.HidePage("dockerPage")
				d.pagesMain.ShowPage("containerLogInfoPage")

				d.app.SetFocus(d.text)

				d.lastFocus = d.app.GetFocus()
				res := d.ContainerInspectId(d.containearTable.GetCell(r, 2).Text)
				d.text.SetBorder(true).
					SetTitle("Meta Data for " + d.containearTable.GetCell(r, 0).Text).
					SetTitleColor(tcell.ColorAqua).
					SetBorderColor(tcell.ColorAqua)
				if res.State.Running {
					fmt.Fprintf(d.text, "[green]Status %s \n", res.State.Status)
				} else {
					fmt.Fprintf(d.text, "[red]Status %s \n", res.State.Status)
				}

				fmt.Fprintf(d.text, "[yellow]Name:[green]   %s [white]\n[yellow]RestartCount:[white] %d \n[yellow]Path:[white]   %s \n[yellow]ID:[white]   %s \n[yellow]Created:[white]   %s \n[yellow]Image:[white]   %s \n[yellow]ResolvConfPath:[white]   %s \n[yellow]HostnamePath:[white]   %s \n[yellow]HostsPath:[white]   %s \n[yellow]LogPath:[white]   %s \n[yellow]Driver:[white]   %s \n[yellow]Platform:[white]   %s \n[yellow]MountLabel:[white]   %s \n[yellow]ProcessLabel:[white]   %s \n[yellow]AppArmorProfile:[white]   %s \n",
					res.ContainerJSONBase.Name,
					res.ContainerJSONBase.RestartCount,
					res.ContainerJSONBase.Path,

					res.ContainerJSONBase.ID,
					res.ContainerJSONBase.Created,

					res.ContainerJSONBase.Image,
					res.ContainerJSONBase.ResolvConfPath,
					res.ContainerJSONBase.HostnamePath,
					res.ContainerJSONBase.HostsPath,
					res.ContainerJSONBase.LogPath,
					res.ContainerJSONBase.Driver,
					res.ContainerJSONBase.Platform,
					res.ContainerJSONBase.MountLabel,
					res.ContainerJSONBase.ProcessLabel,
					res.ContainerJSONBase.AppArmorProfile,
				)

				fmt.Fprintln(d.text, "[green]Args[white]")
				for _, v := range res.ContainerJSONBase.Args {
					fmt.Fprintln(d.text, v)
				}

				fmt.Fprintln(d.text, "[green]Mounts[white]")
				for _, v := range res.Mounts {
					fmt.Fprintln(d.text, v)
				}
				fmt.Fprintln(d.text, "[green]Config[white]")
				fmt.Fprintf(d.text, "HostName: %s \nWorkingDir %s \nDomainName: %s \nUser: %s \n[yellow]Image:[white] %s \n",
					res.Config.Hostname,
					res.Config.WorkingDir,
					res.Config.Domainname,
					res.Config.User,
					res.Config.Image)

				fmt.Fprintln(d.text, "[green]Env[white]")
				for _, v := range res.Config.Env {
					fmt.Fprintln(d.text, v)
				}

				return nil

			}
			if i == 2 {
				d.text.Clear()
				d.text.SetBackgroundColor(tcell.Color16)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.pagesMain.HidePage("dockerPage")
				d.pagesMain.ShowPage("containerLogInfoPage")
				d.app.SetFocus(d.text)
				d.text.SetBorder(true).
					SetTitle("Logs for " + d.containearTable.GetCell(r, 0).Text).
					SetTitleColor(tcell.ColorAqua).
					SetBorderColor(tcell.ColorAqua)
				r, err := d.dockerClient.ContainerLogs(context.Background(), d.containearTable.GetCell(r, 2).Text, types.ContainerLogsOptions{ShowStdout: true, Tail: "200"})

				if err != nil {
					panic(err)
				}

				buf := new(strings.Builder)
				io.Copy(buf, r)
				r.Close()
				fmt.Fprintf(d.text, "%s", buf.String())
				return nil
			}
			if i == 3 {

				d.StopContainerById(d.containearTable.GetCell(r, 2).Text)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.containearTable.Clear()
				d.app.SetFocus(d.lastFocus)

				return nil

			}
			if i == 4 {

				d.StartContainerById(d.containearTable.GetCell(r, 2).Text)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.containearTable.Clear()
				d.app.SetFocus(d.lastFocus)

				return nil

			}

		}
		return event
	})

}

func (d *DockerApi) DropDownImageListNavigation() {
	d.dropdownImageList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			// if d.lastFocus == nil {
			// 	d.lastFocus = d.list

			// }
			d.dropdownImageList.SetLabel("Select an option: ")
			d.dropdownImageList.SetCurrentOption(-1)
			d.pagesMain.HidePage("containerLogInfoPage")
			// d.lastFocus = d.containearTable
			d.app.SetFocus(d.lastFocus)

			return nil

		}
		if event.Key() == tcell.KeyEnter {

			i, _ := d.dropdownImageList.GetCurrentOption()
			r, _ := d.containearTableImage.GetSelection()
			// d.lastFocus = d.app.GetFocus()
			if i == 0 {

				d.RemoveImageByID(d.containearTableImage.GetCell(r, 2).Text)

				d.dropdownImageList.SetLabel("Select an option: ")
				d.dropdownImageList.SetCurrentOption(-1)
				d.app.SetFocus(d.lastFocus)

				return nil
			}

		}
		return event
	})

}
