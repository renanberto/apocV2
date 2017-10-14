package history

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//HTMLOutagesHandler returns a html page with all outage history in a given period
func HTMLOutagesHandler(c *gin.Context){
	c.HTML(http.StatusOK, "history_outage.html", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}