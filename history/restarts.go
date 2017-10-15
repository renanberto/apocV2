package history

import (
	"github.com/renanberto/apocV2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const restatsHistory = "restarts_history"

type ClientRestart struct {
	ClientId      string `json:"client_id" binding:"Required"`
	ServerAddress string `json:"server_address"`
	CreationDate  string
	User          string `json:"user" binding:"Required"`
	Description   string `json:"description"`
}

func InputRestartHandler(c *gin.Context) {

	var tmpClientRestart ClientRestart

	c.BindJSON(&tmpClientRestart)
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	tmpClientRestart.CreationDate = utils.GetCurrentDate()
	Mongoconn := session.DB("").C(restatsHistory)
	err := Mongoconn.Insert(tmpClientRestart)

	if err != nil {
		c.String(http.StatusInternalServerError,"Couldn't save the Version :(")
		return
	}

	c.String(http.StatusOK, "Ok, version saved")
}

func HTMLRestartsHandler(c *gin.Context) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(restatsHistory)
	result := []ClientRestart{}
	err := conn.Find(nil).All(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}

	c.HTML(200, "history_restarts.html", gin.H{
		"Restarts": &result,
	})
}