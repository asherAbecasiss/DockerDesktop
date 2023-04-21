package main

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
)

func (d *DockerApi) GetDockerContainer() []types.Container {
	containers, err := d.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	// for _, container := range containers {
	// 	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	// }

	// fmt.Println(containers)
	// f, _ := os.Create("s.json")

	// data, _ := json.Marshal(containers)
	// f.Write(data)
	// f.Close()

	return containers
}

func (d *DockerApi) RestartContainerID(id string) {

	opt := container.StopOptions{}

	err := d.dockerClient.ContainerRestart(context.TODO(), id, opt)

	if err != nil {
		fmt.Println(err)
	}

}

func (d *DockerApi) ContainerInspectId(id string) types.ContainerJSON {

	res, err := d.dockerClient.ContainerInspect(context.TODO(), id)

	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(res)
	return res

}






//Swarm

func (d *DockerApi) GetSwarmNode() []swarm.Node {

	res, err := d.dockerClient.NodeList(context.Background(), types.NodeListOptions{})

	if err != nil {
		fmt.Println(err)
	}
	// log.Println(res)
	return res

}

func (d *DockerApi) GetDockerServices() []swarm.Service {

	res, err := d.dockerClient.ServiceList(context.Background(), types.ServiceListOptions{})

	if err != nil {
		fmt.Println(err)
	}
	// log.Println(res)
	return res

}

func (d *DockerApi) GetServicesLogs(id string) io.ReadCloser {

	res, err := d.dockerClient.ServiceLogs(context.Background(), id, types.ContainerLogsOptions{ShowStdout: true})

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(data)
	// log.Println(res)
	return res

}
func (d *DockerApi) DockerServicesUpdate(id string) {

	// inspect, _, err := d.dockerClient.ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions{})
	// fmt.Println(inspect.ID)
	// fmt.Println("------")

	// // obj := Msg{f}
	// // json, _ := json.Marshal(obj)
	// // fmt.Println(string(json))
	// // myString := string(f[:])
	// // fmt.Println(myString)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // s := inspect.PreviousSpec
	// // fmt.Println(*s)

	// res, err := d.dockerClient.ServiceUpdate(context.Background(), inspect.ID, swarm.Version{Index: inspect.Meta.Version.Index}, inspect.Spec, types.ServiceUpdateOptions{RegistryAuthFrom: "previous-spec"})

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("------")

	d.dockerClient.ServiceRemove(context.Background(), id)
	// log.Println(res.Warnings)

}
