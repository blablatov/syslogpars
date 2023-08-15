package syslog2mongo

import (
	"fmt"
	"log"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	//"github.com/globalsign/mgo/dbtest"
	//"gotest.tools/v3/assert"
)

func TestStrMongo(t *testing.T) {
	var tests = []struct {
		AddrData  string
		AlarmData string
	}{
		{"addr_1", "maximum temp"},
		{"addr_2", "min temp"},
		{"addr_54", "max temp"},
		{"Data for test", "norm temp"},
		{"Yes, true", "high temp"},
	}

	var prevAddrData string
	for _, test := range tests {
		if test.AddrData != prevAddrData {
			fmt.Printf("\n%s\n", test.AddrData)
			prevAddrData = test.AddrData
		}
	}

	var prevAlarmData string
	for _, test := range tests {
		if test.AlarmData != prevAlarmData {
			fmt.Printf("\n%s\n", test.AlarmData)
			prevAlarmData = test.AlarmData
		}
	}
}

func TestSendMongo(t *testing.T) {
	DsnMongo := "mongodb://localhost:27017/syslogdb"
	session, err := mgo.Dial(DsnMongo)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// is check name in dBase
	c := session.DB("syslogdb").C("syslogdata")
	chk := DataMongo{}
	err = c.Find(bson.M{"sensortype": chk.AddrData, "datasensor": chk.AlarmData}).One(&chk)
	if err != nil {
		log.Print("\nErr data for write to MongoDB, method of interface: ", err)
		err = c.Insert(&DataMongo{chk.AddrData, chk.AlarmData})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sensor was written via method of interface:", chk.AddrData,
			"\nData of sensor was written via method of interface:", chk.AlarmData)
	}
}

// // DB Server controls a MongoDB server process to used within test suites
// type M map[string]interface{}

// func TestSessionMongo(t *testing.T) {
// 	var server dbtest.DBServer
// 	server.SetPath("./")
// 	defer server.Stop()

// 	session := server.Session()
// 	err := session.DB("mongodb").C("mongodata").Insert(M{"m": 1})
// 	session.Close()
// 	assert.Assert(t, err != nil)

// 	server.Wipe()

// 	session = server.Session()
// 	names, err := session.DatabaseNames()
// 	session.Close()
// 	assert.Assert(t, err != nil)
// 	for _, name := range names {
// 		if name != "local" && name != "admin" {
// 			log.Fatalf("Wipe should have removed this database: %s", name)
// 		}
// 	}
// }
