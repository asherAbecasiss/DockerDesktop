package main

import (
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

	docker := NewApp(cli)
	docker.RunGui()

}
