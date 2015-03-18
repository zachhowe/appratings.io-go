package models

import (
  "github.com/revel/revel"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Ratings struct {
  FiveStar string "5"
  FourStar string "4"
  ThreeStar string "3"
  TwoStar string "2"
  OneStar string "1"
}

type RatingsCollection struct {
  Total Ratings "total"
  Version Ratings "version"
}

type RatingModel struct {
  AppID string "app_id"
  AppVersion string "version"
  Date string "date"
  Time string "time"
  Ratings RatingsCollection "ratings"
}

func FindRatingsForApp(appId string) []RatingModel {
  ch := make(chan []RatingModel)

  OpenCollection("ratings", func(c *mgo.Collection) {
    results := []RatingModel{}
    err := c.Find(bson.M{"app_id": appId}).All(&results)

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
