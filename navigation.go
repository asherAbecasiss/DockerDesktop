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
			d.lastFocus = d.app.GetFocus()
			
			d.pagesMain.SwitchToPage("mainPage")
			d.app.SetFocus(d.list)
			return event
		}
		return event
	})
	d.text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {

			d.pagesTabel.SwitchToPage("containerTablePage")
			d.app.SetFocus(d.table)

		}
		return event
	})
	d.dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyESC {
			// if d.lastFocus == nil {
			// 	d.lastFocus = d.list

			// }
			d.dropdown.SetLabel("Select an option: ")
			d.dropdown.SetCurrentOption(-1)
			d.pagesTabel.ShowPage("containerTablePage")
			d.lastFocus = d.table
			d.app.SetFocus(d.table)
		}
		if event.Key() == tcell.KeyEnter {

			i, _ := d.dropdown.GetCurrentOption()
			r, _ := d.table.GetSelection()
			d.lastFocus = d.app.GetFocus()
			if i == 0 {

				d.RestartContainerID(d.table.GetCell(r, 2).Text)

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.app.SetFocus(d.lastFocus)

				return nil
			}
			if i == 1 {
				d.text.Clear()

				d.dropdown.SetLabel("Select an option: ")
				d.dropdown.SetCurrentOption(-1)
				d.pagesTabel.ShowPage("containerLogInfoPage")

				d.app.SetFocus(d.text)
				res := d.ContainerInspectId(d.table.GetCell(r, 2).Text)
				d.text.SetBorder(true).
					SetTitle("Meta Data for " + d.table.GetCell(r, 0).Text).
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
				d.pagesTabel.ShowPage("containerLogInfoPage")
				d.app.SetFocus(d.text)
				d.text.SetBorder(true).
					SetTitle("Logs for " + d.table.GetCell(r, 0).Text).
					SetTitleColor(tcell.ColorAqua).
					SetBorderColor(tcell.ColorAqua)
				r, err := d.dockerClient.ContainerLogs(context.Background(), d.table.GetCell(r, 2).Text, types.ContainerLogsOptions{ShowStdout: true, Tail: "200"})

				if err != nil {
					panic(err)
				}

				buf := new(strings.Builder)
				io.Copy(buf, r)
				r.Close()
				fmt.Fprintf(d.text, "%s", buf.String())
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
			// d.pagesTabel.ShowPage("main1")
			d.lastFocus = d.app.GetFocus()
			d.app.SetFocus(d.lastFocus)

		}

		return event
	})

}
