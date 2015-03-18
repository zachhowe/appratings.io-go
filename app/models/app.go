package models

import (
  "github.com/revel/revel"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type AppInfo struct {
  TrackName string "trackName"
}

type AppModel struct {
  AppID string "app_id"
  Info AppInfo "info"
}

func FindAllApps() []AppModel {
  ch := make(chan []AppModel)

  OpenCollection("apps", func(c *mgo.Collection) {
    results := []AppModel{}
    err := c.Find(bson.M{}).All(&results)

    if err != nil {
      revel.ERROR.Printf("Error finding apps: %s", err)
    } else {
      revel.INFO.Printf("Results: %s", results)
    }

    ch <- results
  })

  return <-ch
}

func init() {
  revel.OnAppStart(func() {
  })
}
