package history

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/renanberto/apocV2/utils"
)

const tickets_collection = "open_tickets"

type Ticket struct {
	ClientID     string
	CreationDate string
	Open         bool
}

func New(clientId, creationDate string, open bool) Ticket {
	ticket := Ticket{clientId, creationDate, open}
	return ticket
}

func GetOpenTickets(clientId string) (Ticket, error) {
	var ticket Ticket
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(tickets_collection)
	err := conn.Find(bson.M{"clientid": clientId, "open": true}).One(&ticket)
	if err != nil {
		return Ticket{}, fmt.Errorf("Something went wrong error: %e", err)
	}
	return ticket, err
}

func (t *Ticket) OpenTicket(clientId string) error {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(tickets_collection)
	err := conn.Insert(t)
	return err
}

func (t *Ticket) CloseTicket() error {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C(tickets_collection)

	query := bson.M{"clientid": t.ClientID, "open": true}
	change := bson.M{"$set": bson.M{"open": false}}

	err := conn.Update(query, change)

	return err
}

func (t *Ticket) Validate() bool {
	if t.Open == false {
		return false
	}
	return true
}

func HTMLTicketsHandler(c *gin.Context) {
	session := utils.Mongoconnect().Copy()
	defer session.Close()

	conn := session.DB("").C("open_tickets")

	result := []Ticket{}

	err := conn.Find(nil).All(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
	}

	c.HTML(http.StatusOK, "history_tickets.html", gin.H{
		"Tickets": &result,
		"title": "Agile Promoter Operations Center",
	})
}