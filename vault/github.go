package vault

import (
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"time"
)

var (
	githubEndpoint="/v1/auth/github/login"
	vaultAddr = os.Getenv("VAULT_ADDR")
)

type PolicyAuthGithubResponse map[string]string

type GithubResponse struct {
	RequestId				string							`json:"request_id"`
	Auth						AuthGithubResponse	`json:"auth"`
}

type AuthGithubResponse struct {
	ClientToken			string											`json:"client_token"`
	Accessor				string											`json:"accessor"`
	Policies				PolicyAuthGithubResponse 		`json:"policies"`
	Metadata				MetadataAuthGithubResponse
	LeaseDuration		int													`json:"lease_duration"`
	Renewable				bool
}

type MetadataAuthGithubResponse struct {
	Org				string	`json:"org"`
	Username	string	`json:"username"`
}

func GithubLogin(githubToken string) GithubResponse {

	var tmpGithubResponse GithubResponse

	githubUrlLogin := vaultAddr + githubEndpoint
	value := map[string]string{"token": githubToken}
	githubTokenJson, _ := json.Marshal(value)

	req, _ := http.NewRequest("POST", githubUrlLogin, bytes.NewBuffer(githubTokenJson))

	client := &http.Client{Timeout: time.Second * 10}
	resp, _ := client.Do(req)
	bodyGithubResponse , _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(bodyGithubResponse), &tmpGithubResponse)

	return tmpGithubResponse
}