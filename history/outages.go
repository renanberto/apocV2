package history

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/renanberto/apocV2/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

//HTMLOutagesHandler returns a html page with all outage history in a given period
const okState = "OK"
const criticalState = "CRITICAL"
const outageHistory = "outage_history"
const outageCounter = "outage_counter"

//OutageRecord is the structure for an outage message provenient from Nagios or another monitoring tool
type OutageRecord struct {
	ClientID      string `json:"client_id" binding:"Required"`
	CreationDate  string
	ServerAddress string `json:"server_address"`
	Description   string `json:"description"`
	TopClients    string `json:"top_clients"`
	User          string `json:"user" binding:"Required"`
	Agent         string `json:"agent" binding:"Required"`
	State         string `json:"state" binding:"Required"`
	Duration      time.Duration
}

//InputOutageHandler saves an outage message to the outage_history collection
func InputOutageHandler(c *gin.Context) {
	var client OutageRecord

	c.BindJSON(&client)
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	client.CreationDate = utils.GetCurrentDate()
	conn := session.DB("").C(outageHistory)
	err := conn.Insert(client)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't save outage history")
		return
	}

	tkts, err := GetOpenTickets(client.ClientID)
	checkError(err)

	if (tkts.Validate()) && (client.State == okState) {
		err := saveOutage(client, tkts)
		checkError(err)
		err = tkts.CloseTicket()
		checkError(err)
		c.String(http.StatusOK, "Outage Saved!")
		return
	} else if (tkts.Validate()) && (client.State == criticalState) {
		c.String(http.StatusOK, "There's already a open ticket for the client "+client.ClientID)
		return
	} else if !tkts.Validate() && (client.State == okState) {
		c.String(http.StatusOK, "This message with ok state should not arrived here")
		return
	}

	ticket := New(client.ClientID, utils.GetCurrentDate(), true)
	err = ticket.OpenTicket(client.ClientID)
	checkError(err)

	c.String(http.StatusOK, "Ticket Opened for client "+client.ClientID)
	return
}

//HTMLOutagesHandler returns a html page with all outage historys in a given period
func HTMLOutagesHandler(c *gin.Context) {
	result := []OutageRecord{}
	var err error

	session := utils.Mongoconnect().Clone()
	defer session.Close()

	beginDate := c.Query("beginDate")
	endDate := c.Query("endDate")

	if endDate != "" && beginDate != "" {
		endDate += "T23:59:00Z"
		beginDate += "T00:00:00Z"
		result, err = getOutagesByDate(beginDate, endDate)
		checkError(err)
	} else {
		result, err = getRecentOutages()
		checkError(err)
	}

	c.HTML(http.StatusOK, "history_outage.html", gin.H{
		"Outages": &result,
		"title":   "Agile Promoter Operations Center",
	})
}

//saveOutage saves an outage message to the outage_counter collection when a "OK" message comes from Nagios
func saveOutage(client OutageRecord, ticket Ticket) error {

	startTime, err := time.Parse(time.RFC3339, ticket.CreationDate)
	checkError(err)
	endTime, err := time.Parse(time.RFC3339, utils.GetCurrentDate())
	checkError(err)

	client.Duration = endTime.Sub(startTime)

	conn := utils.Mongoconnect().DB("").C(outageCounter)
	err = conn.Insert(client)
	return err
}

// getRecentOutages search log in mongoDB
func getRecentOutages() ([]OutageRecord, error) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	result := []OutageRecord{}

	aDay := 24 * time.Hour
	beginofToday := time.Now().UTC().Truncate(aDay)
	today := beginofToday.Add(23 * time.Hour).Format(time.RFC3339)

	beginOflastweek := time.Now().UTC().Add(-168 * time.Hour)
	beginOflastweek = beginOflastweek.Truncate(aDay)

	lastWeek := beginOflastweek.Format(time.RFC3339)

	//limite no mongodb
	conn := session.DB("").C(outageCounter)
	err := conn.Find(
		bson.M{"creationdate": bson.M{
			"$gte": lastWeek,
			"$lt":  today,
		},
		},
	).All(&result)

	return result, err
}

// getOutagesByDate search log in mongoDB by date
func getOutagesByDate(begin, end string) ([]OutageRecord, error) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	outages := []OutageRecord{}
	conn := session.DB("").C(outageCounter)
	err := conn.Find(
		bson.M{"creationdate": bson.M{
			"$gte": begin,
			"$lt":  end,
		},
		},
	).All(&outages)

	return outages, err
}

// checkError is handler error
func checkError(err error) {
	if err != nil {
		fmt.Errorf("Something went wrong, the error ir %e", err)
	}
}
