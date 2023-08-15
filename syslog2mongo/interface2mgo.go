package syslog2mongo

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ApcMonger interface {
	SendMongo(DsnMongo string) (bool, error)
}

type ApcMongo struct {
	AddrData  string
	AlarmData string
}

func (md ApcMongo) SendMongo(DsnMongo string) (bool, error) {
	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		log.Fatalf("Error conn to Mongo: %v", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("syslogdb").C("syslogdata")
	chk := ApcMongo{}
	err = c.Find(bson.M{"addr": md.AddrData, "alarmdata": md.AlarmData}).One(&chk)
	if err == nil {
		log.Println("\nData already was to DB via method of interface", err)
		return false, err
	}
	if err != nil {
		log.Print("\nErr data for write to MongoDB, method of interface: ", err)
		err = c.Insert(&ApcMongo{md.AddrData, md.AlarmData})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nSensor was written via method of interface:", md.AddrData,
			"\nData of sensor was written via method of interface:", md.AlarmData)
	}
	return true, err
}
