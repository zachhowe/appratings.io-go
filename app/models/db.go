package models

import (
  "github.com/revel/revel"
  "gopkg.in/mgo.v2"
)

type OpenCollectionHandler func(*mgo.Collection)

func OpenCollection(collectionName string, handler OpenCollectionHandler) {
  // fetch config values
  db_host := revel.Config.StringDefault("db.host", "localhost")
  db_name := revel.Config.StringDefault("db.name", "appratings")

  // start
  session, err := mgo.Dial(db_host)
  if err != nil {
    panic(err)
  }
  session.SetMode(mgo.Monotonic, true)
  c := session.DB(db_name).C(collectionName)

  go func(session *mgo.Session) {
    defer session.Close()
    
    handler(c)
  }(session)
}
