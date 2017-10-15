package history

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/renanberto/apocV2/utils"
)

const updateHistory = "update_history"

type ClientUpdate struct {
	ClientId        string `json:"client_id" binding:"Required"`
	Platform        string `json:"platform" binding:"Required"`
	Date            string `json:"date"`
	PreviousVersion string `json:"previous_version" binding:"Required"`
	UpdatedVersion  string `json:"updated_version" binding:"Required"`
	User            string `json:"user" binding:"Required"`
}

func InputUpdateHandler(c *gin.Context) {

	var tmpClientUpdate ClientUpdate

	c.BindJSON(&tmpClientUpdate)
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(updateHistory)
	err := conn.Insert(tmpClientUpdate)

	if err != nil {
		c.String(http.StatusInternalServerError,"Couldn't save update history :(")
		return
	}
	c.String(http.StatusOK,"Update history Saved")
	return
}

func InputUpdateHandlerByID(c *gin.Context) {
	result := []ClientUpdate{}

	session := utils.Mongoconnect().Copy()
	defer session.Close()
	conn := session.DB("").C(updateHistory)

	err := conn.Find(bson.M{"id": c.Param(":client_id")}).All(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}
	c.JSON(http.StatusOK, &result)
}

func HTMLUpdatesHandler(c *gin.Context) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(updateHistory)

	result := []ClientUpdate{}

	err := conn.Find(nil).All(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}
	c.HTML(http.StatusOK,"history_updates.html", gin.H{
		"Updates": &result,
	})
}