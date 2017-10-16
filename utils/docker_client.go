package utils

import (
	"os"
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

var (
	path = os.Getenv("DOCKER_CERT_PATH")
)

type DockerConnection struct {
	endpoint  string
	ca 				string
	cert   		string
	key  			string
}

func (c *DockerConnection) New(clientServerName string) {
	c.endpoint = "tcp://" + clientServerName + ".agilepromoter.com:2376"
	c.ca = fmt.Sprintf("%s/ca.pem", path)
	c.cert = fmt.Sprintf("%s/cert.pem", path)
	c.key = fmt.Sprintf("%s/key.pem", path)
}

func (c *DockerConnection) GetProdutionContainerId(dockerConnection DockerConnection, clientId string) (containerId string) {
	client, _ := docker.NewTLSClient(dockerConnection.endpoint,dockerConnection.cert,dockerConnection.key,dockerConnection.ca)
	allContainersRunning, _ := client.ListContainers(docker.ListContainersOptions{All: false})

	clientId = "/production_" + clientId

	for _, element := range allContainersRunning {
		for _, names := range element.Names {
			if clientId == names {
				containerId := element.ID
				return containerId
			}
		}
	}
	return containerId
}

func RestartContainerId (containerId string, dockerConnection DockerConnection) (err error){
	client, _ := docker.NewTLSClient(dockerConnection.endpoint,dockerConnection.cert,dockerConnection.key,dockerConnection.ca)
	err = client.RestartContainer(containerId,10)

	return err
}

func RemoveContainerId (containerId string, dockerConnection DockerConnection ) (err error){
	client, _ := docker.NewTLSClient(dockerConnection.endpoint,dockerConnection.cert,dockerConnection.key,dockerConnection.ca)
	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: containerId, RemoveVolumes: true, Force: true})

	return err
}

func ListAllContainers (dockerConnection DockerConnection) (err error, listContainers []docker.APIContainers) {
	client, _ := docker.NewTLSClient(dockerConnection.endpoint, dockerConnection.cert, dockerConnection.key, dockerConnection.ca)
	allContainersRunning, _ := client.ListContainers(docker.ListContainersOptions{All: false})

	return err, allContainersRunning
}