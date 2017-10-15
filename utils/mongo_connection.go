package utils

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	mongoAddr     = os.Getenv("MONGO_DB_ADDR")
	mongoUser     = os.Getenv("MONGO_DB_USER")
	mongoPassword = os.Getenv("MONGO_DB_PASSWORD")
	mongoDatabase = os.Getenv("MONGO_DB_NAME")
)

//Mongoconnect will connect to a database and return it's session
func Mongoconnect() *mgo.Session {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddr},
		Timeout:  5 * time.Second,
		Database: mongoDatabase,
		Username: mongoUser,
		Password: mongoPassword,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal("Couldn't connect to mongodb :( ", err)
	}
	return session
}

func Connect(c *gin.Context) {
	s := Mongoconnect().Clone()

	defer s.Close()
	c.Set("db", s.DB(mongoDatabase))
	c.Next()
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.HTML(http.StatusBadRequest, "400", gin.H{
			"errors": c.Errors,
		})
	}
}
