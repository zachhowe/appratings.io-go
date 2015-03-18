package models

import (
  "github.com/revel/revel"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type UserModel struct {
  UserName string "username"
  Password string "password"
}

func AuthenticateUser(username string, password string) []UserModel {
  ch := make(chan []UserModel)

  OpenCollection("users", func(c *mgo.Collection) {
    results := []UserModel{}
    err := c.Find(bson.M{}).All(&results)

    if err != nil {
      revel.ERROR.Printf("Error finding users: %s", err)
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
