package staging

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/renanberto/apocV2/utils"
	"strings"
	"time"
	"fmt"
)

type ContainerInformation struct {
	ID	string
	Name	string
	Status	string
	State	string
	Created	time.Time
}

func HTMLStagingHandler(c *gin.Context) {

	listedStagingContainers := listAllStagingContainers()

	c.HTML(http.StatusOK, "staging.html", gin.H{
		"title": "Agile Promoter Operations Center",
		"stagingContainers": &listedStagingContainers,
	})
}

func listAllStagingContainers() []ContainerInformation {
	  var newDockerConnection utils.DockerConnection
		var stagingContainersInformation []ContainerInformation
		var ci ContainerInformation

		newDockerConnection.New("involves-hetfield")

		err, stagingContainers := utils.ListAllContainers(newDockerConnection)
		if err != nil {panic(err)}

		for _, element := range stagingContainers {
			joinContainerName := strings.Join(element.Names,"")
			ContainerNameWithReplace := strings.Replace(joinContainerName,"/","",1)
			t := time.Unix(element.Created,0)
			if strings.Contains(ContainerNameWithReplace,"staging_qa_") {
				ci.ID = element.ID
				ci.Name = ContainerNameWithReplace
				ci.Status = element.Status
				ci.Created = t
				ci.State = element.State
				stagingContainersInformation = append(stagingContainersInformation,ci)
			}
		}

	return stagingContainersInformation
}

func RemoveStagingContainer(c *gin.Context) {
	var newDockerConnection utils.DockerConnection

	containerID := c.PostForm("id")
	newDockerConnection.New("involves-hetfield")
	err := utils.RemoveContainerId(containerID,newDockerConnection)

	fmt.Println("passou")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": "success",
		})
	}

}