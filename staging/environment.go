package staging

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/renanberto/apocV2/utils"
	"fmt"
)

func HTMLStagingHandler(c *gin.Context) {

	listedStagingContainers := listAllStagingContainers()
	for _,i := range listedStagingContainers {

		fmt.Println(i)
	}

	c.HTML(http.StatusOK, "staging.html", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}

func listAllStagingContainers() (listedStagingContainers []string){
	  var tmpDockerConnection utils.DockerConnection
	  var stagingNames []string

	  tmpDockerConnection.New("involves-hetfield")

		err, allStagingContainers := utils.ListAllContainers(tmpDockerConnection)
		if err != nil {panic(err)}

		for _, element := range allStagingContainers {
			for _, name := range element.Names {
				stagingNames = append(stagingNames,name)
			}
		}

	return stagingNames
}