package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// Ipv4 is the main record format
type Ipv4 struct {
	//	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	IP          string       `bson:"ip" json:"ip"`
	Assignments []Assignment `bson:"assignments" json:"assignments"`
}

// Assignment contains each individual assignment date/time
type Assignment struct {
	Mac  string    `bson:"mac" json:"mac"`
	Time time.Time `bson:"time" json:"time"`
}

func addAssignment(ipv4 string, assignment Assignment, s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("dhcplogger").C("ips")

	pushToArray := bson.M{"$addToSet": bson.M{"assignments": assignment}}
	q := bson.M{"ip": ipv4}
	_, err := c.Upsert(q, pushToArray)
	if err != nil {
		log.Println("Error, Problem adding Assignment to Database: ", err)
	} else {
		log.Println("DHCPv4," + assignment.Mac + "," + ipv4)
	}

}
