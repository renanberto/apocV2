package main

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "github.com/involvestecnologia/apoc/v2/history"
  "github.com/involvestecnologia/apoc/v2/vault"
)

func main() {
  router := gin.Default()
  router.LoadHTMLGlob("templates/*")
  router.Static("/dist", "public/dist")

  // Posts
  router.POST("/vault/mysql/create-creds", vault.InputMysqlHandler)

  // Views
  router.GET("/", HTMLIndexHandler)
  router.GET("/about", HTMLAboutHandler)
  router.GET("/history/outages", history.HTMLOutagesHandler)
  router.GET("/vault/mysql", vault.HTMLMysqlHandler)

  router.Run(":4000")
}

func HTMLIndexHandler(c *gin.Context) {
  c.HTML(http.StatusOK, "index.tmpl", gin.H{
    "title": "Agile Promoter Operations Center",
  })
}

func HTMLAboutHandler(c *gin.Context) {
  c.HTML(http.StatusOK, "about.tmpl", gin.H{
    "title": "Agile Promoter Operations Center",
  })
}
