package history

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/renanberto/apocV2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const versionsHistory = "versions_control"

type Versions struct {
	WebVersion     string `json:"web_version" binding:"Required"`
	IosVersion     string `json:"ios_version" binding:"Required"`
	AndroidVersion string `json:"android_version" binding:"Required"`
}

func HTMLVersionsHandler(c *gin.Context) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(versionsHistory)
	result := []Versions{}

	err := conn.Find(nil).All(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}

	c.HTML(http.StatusOK, "history_versions.html", gin.H{
		"Versions": &result,
	})
}

func InputVersionsHandler(c *gin.Context) {

	var tmpVersions Versions

	c.BindJSON(&tmpVersions)
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	Mongoconn := session.DB("").C(versionsHistory)
	err := Mongoconn.Insert(tmpVersions)

	if err != nil {
		c.String(http.StatusInternalServerError,"Couldn't save the Version :(")
		return
	}
	c.String(http.StatusOK,"Ok, version saved")
	return
}

func InputUpdateVersionsHandler(c *gin.Context) {

	var tmpVersion Versions

	c.BindJSON(&tmpVersion)
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	Mongoconn := session.DB("").C(versionsHistory)
	err := Mongoconn.Update(bson.M{"webversion": c.Request.Form.Get("old_web_version")}, tmpVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal error")
	} else {
		c.Redirect(http.StatusOK,"/history/versions")
	}

}

func InputVersionHandlerByWebVersion(c *gin.Context) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	result := []Versions{}
	Mongoconn := session.DB("").C(versionsHistory)

	err := Mongoconn.Find(bson.M{"webversion": c.Param("webversion")}).All(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}

	c.JSON(http.StatusOK, &result)

}
