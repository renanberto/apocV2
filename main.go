package main

import (
	"github.com/gin-gonic/gin"
	"github.com/renanberto/apocV2/ginHtmlRender"
	"github.com/renanberto/apocV2/history"
	"github.com/renanberto/apocV2/utils"
	"github.com/renanberto/apocV2/vault"
	"net/http"
)

func main() {
	utils.Mongoconnect()

	router := gin.Default()
	router.Use(utils.Connect)
	router.Use(utils.ErrorHandler)

	htmlSettings(router)

	// Posts
	v1 := router.Group("/v1")
	{
		endPointHistory := router.Group("/history")
		{

			endPointHistory.POST("/outage/create", history.InputOutageHandler)
			endPointHistory.POST("/restart/create", history.InputRestartHandler)
			endPointHistory.POST("/update/create", history.InputUpdateHandler)
			endPointHistory.POST("/versions/update", history.InputUpdateVersionsHandler)
			endPointHistory.POST("/versions", history.InputVersionsHandler)
		}
		v1.POST("/vault/mysql/create-creds", vault.InputMysqlHandler)
	}

	// Views
	router.GET("/", HTMLIndexHandler)
	router.GET("/about", HTMLAboutHandler)
	router.GET("/history/outages", history.HTMLOutagesHandler)
	router.GET("/history/tickets", history.HTMLTicketsHandler)
	router.GET("/history/restarts", history.HTMLRestartsHandler)
	router.GET("/history/updates", history.HTMLUpdatesHandler)
	router.GET("/history/versions", history.HTMLVersionsHandler)
	router.GET("/vault/mysql", vault.HTMLMysqlHandler)

	//With Params
	router.GET("/v1/history/update/:client_id", history.InputUpdateHandlerByID)
	router.GET("/v1/history/versions/:webversion", history.InputVersionHandlerByWebVersion)

	router.Run(":4000")
}

func htmlSettings(router *gin.Engine) {
	htmlRender := ginHtmlRender.New()
	htmlRender.Debug = gin.IsDebugging()
	htmlRender.Layout = "index"

	router.HTMLRender = htmlRender.Create()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true

	router.LoadHTMLGlob("templates/*")
	router.Static("/dist", "public/dist")
}

func HTMLIndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}

func HTMLAboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}
