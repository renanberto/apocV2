package vault

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var (
	MongoEndpoint = "/v1/database/creds/"
)

type MongoResponse struct {
	LeaseDuration int       `json:"lease_duration"`
	Data          MysqlData `json:"data"`
	Errors        []string  `json:"errors"`
}

type MongoData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MongoCredsForm struct {
	githubToken string `json:"githubToken" binding:"required"`
	database    string `json:"database" binding:"required"`
	accessMode  string `json:"accessMode" binding:"required"`
}

// This function HTMLMysqlHandler provider a HTML to MYSQL creds generator
func HTMLMongoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "vault_mongo.html", gin.H{
		"title": "Agile Promoter Operations Center",
	})
}

// This function InputMysqlHandler call the POST action
func InputMongoHandler(c *gin.Context) {

	var tmpMongoCredsForm MongoCredsForm

	tmpMongoCredsForm.githubToken = c.PostForm("githubToken")
	tmpMongoCredsForm.database = c.PostForm("database")
	tmpMongoCredsForm.accessMode = c.PostForm("accessMode")

	githubBody := GithubLogin(tmpMongoCredsForm.githubToken)
	mongoCredsInformation := GenerateMongoCreds(githubBody.Auth.ClientToken, tmpMongoCredsForm.database, tmpMongoCredsForm.accessMode)

	if githubBody.Errors != nil {
		c.JSON(http.StatusOK, gin.H{
			"respError": githubBody.Errors,
		})
		return
	}

	if mongoCredsInformation.Errors != nil {
		c.JSON(http.StatusOK, gin.H{
			"respError": mongoCredsInformation.Errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":   mongoCredsInformation.Data.Username,
		"password":   mongoCredsInformation.Data.Password,
		"lease_time": mongoCredsInformation.LeaseDuration,
	})
}

// This function GenerateMysqlCreds creates a credentials of MYSQL
func GenerateMongoCreds(clientToken, database, accessMode string) MongoResponse {

	var tmpMongoResponse MongoResponse

	githubUrlLogin := vaultAddr + MongoEndpoint + accessMode + "_" + database

	req, _ := http.NewRequest("GET", githubUrlLogin, nil)
	req.Header.Add("X-Vault-Token", clientToken)
	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	bodyMongoResponse, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(bodyMongoResponse), &tmpMongoResponse)

	return tmpMongoResponse
}
