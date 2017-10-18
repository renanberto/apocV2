package utils

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"os"
)

var (
	jenkinsUrl        = os.Getenv("JENKINS_URL")
	jenkinsAuthBase64 = b64.URLEncoding.EncodeToString([]byte(os.Getenv("JENKINS_AUTH")))
)

type ReviewApp struct {
	pull_id_java string
	pull_id_html string
	user         string
}

func (c *ReviewApp) NewReviewApp(pullIdJava, pullIdHtml, user string) {
	c.pull_id_java = pullIdJava
	c.pull_id_html = pullIdHtml
	c.user = user
}

func (c *ReviewApp) PostJenkins(Url string) error {
	req, err := http.NewRequest("POST", Url, nil)
	req.Header.Add("Authorization", "Basic "+jenkinsAuthBase64)
	q := req.URL.Query()

	q.Add("SLACK_USER", c.user)
	q.Add("PULL_ID_HTML", c.pull_id_html)
	q.Add("PULL_ID_JAVA", c.pull_id_java)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		return error(fmt.Errorf("Request failed, Status code: %d", res.StatusCode))
	}

	return err

}

func (c *ReviewApp) ReviewAppCreate() {
	path = jenkinsUrl + "/job/agilepromoter-staging-pr-create/buildWithParameters"
	c.PostJenkins(path)
}