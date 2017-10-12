package vault

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	mysqlEndpoint="/v1/database/creds/"
)

type MysqlResponse struct {
	LeaseDuration		int				`json:"lease_duration"`
	Data						MysqlData `json:"data"`
}

type MysqlData struct {
	Username			string	`json:"username"`
	Password			string	`json:"password"`
}

type MysqlCredsForm struct {
	githubToken		string `json:"githubToken" binding:"required"`
	database 			string `json:"database" binding:"required"`
	accessMode 		string `json:"accessMode" binding:"required"`
}

// This function HTMLMysqlHandler provider a HTML to MYSQL creds generator
func HTMLMysqlHandler(c *gin.Context){
	c.HTML(http.StatusOK, "vault_mysql.tmpl", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}

// This function InputMysqlHandler call the POST action
func InputMysqlHandler(c *gin.Context) {

	var tmpMysqlCredsForm MysqlCredsForm

	tmpMysqlCredsForm.githubToken = c.PostForm("githubToken")
	tmpMysqlCredsForm.database = c.PostForm("database")
	tmpMysqlCredsForm.accessMode = c.PostForm("accessMode")

	githubBody := GithubLogin(tmpMysqlCredsForm.githubToken)
	mysqlCredsInformation := GenerateMysqlCreds(githubBody.Auth.ClientToken, tmpMysqlCredsForm.database, tmpMysqlCredsForm.accessMode)

	username := mysqlCredsInformation.Data.Username
	password := mysqlCredsInformation.Data.Password

	c.JSON(http.StatusOK,gin.H{
		"username": username,
		"password": password,
	})
}

// This function GenerateMysqlCreds creates a credentials of MYSQL
func GenerateMysqlCreds(clientToken, database, accessMode string) MysqlResponse {

	var tmpMysqlResponse MysqlResponse

	githubUrlLogin := vaultAddr + mysqlEndpoint + accessMode + "_" + database

	req, err := http.NewRequest("GET", githubUrlLogin, nil)
	req.Header.Add("X-Vault-Token", clientToken)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyMysqlResponse , _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(bodyMysqlResponse), &tmpMysqlResponse)

	return tmpMysqlResponse
}