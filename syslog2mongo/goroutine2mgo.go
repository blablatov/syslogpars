package syslog2mongo

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataMongo struct {
	AddrData  string
	AlarmData string
}

func SendgMongo(AddrData, DsnMongo, AlarmData string, adr, alr chan string, wg sync.WaitGroup) {
	defer wg.Done()

	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("syslogdb").C("syslogdata")
	chk := DataMongo{}
	err = c.Find(bson.M{"addr": AddrData, "alarmdata": AlarmData}).One(&chk)
	if err == nil {
		log.Println("\nData already was to DB via goroutine", err)
		return
	}
	if err != nil {
		log.Print("\nData in MongoDB via goroutine: ", err)
		err = c.Insert(&DataMongo{AddrData, AlarmData})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nSyslog was written via goroutine:", AddrData, "\nData was written via goroutine:", AlarmData)
		adr <- AddrData
		alr <- AlarmData
	}
}
