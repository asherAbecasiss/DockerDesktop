package main

import (
	"fmt"
	"net"

	"github.com/docker/docker/client"
)

func NewApp(client *client.Client) *DockerApi {
	return &DockerApi{
		dockerClient: client,
	}
}

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	docker := NewApp(cli)
	docker.RunGui()
	// localAddresses()

}
func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				fmt.Printf("-->%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())

			case *net.IPNet:
				fmt.Printf("------->%v : %s llllllll [%v/%v]\n", i.Name, v, v.IP, v.Mask)
			}

		}
	}
}
