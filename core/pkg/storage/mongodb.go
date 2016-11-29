package storage

import (
	"log"

	"gopkg.in/mgo.v2"
)

var hosts="localhost"

type DataStore struct {
    session *mgo.Session
}
func ConstructDataStore() (*DataStore,error){
    session, err := mgo.Dial(hosts)
	if err != nil {
		log.Fatal(err)
        return nil,err
	}
    return &DataStore{session},nil
}

func (ds *DataStore) Insert(collection string, obj interface{}) {
	ds.session.SetMode(mgo.Monotonic, true)

	c := ds.session.DB("leo").C(collection)

    if err := c.Insert(obj);err!=nil{
        log.Println("insert error: ", err)
    }
    
}
