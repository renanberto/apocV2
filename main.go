package main

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "github.com/renanberto/apocV2/history"
  "github.com/renanberto/apocV2/vault"
  "github.com/renanberto/apocV2/utils"
  "github.com/renanberto/apocV2/ginHtmlRender"
 )

func init() {
  utils.Mongoconnect()
}

func main() {
  router := gin.Default()
  router.Use(utils.Connect)
  router.Use(utils.ErrorHandler)

  htmlSettings(router)

  // Posts
  router.POST("/vault/mysql/create-creds", vault.InputMysqlHandler)

  // Views
  router.GET("/", HTMLIndexHandler)
  router.GET("/about", HTMLAboutHandler)
  router.GET("/history/outages", history.HTMLOutagesHandler)
  router.GET("/vault/mysql", vault.HTMLMysqlHandler)

  router.Run(":4000")
}

func htmlSettings(router *gin.Engine){
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
