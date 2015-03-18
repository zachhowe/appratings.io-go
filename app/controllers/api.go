package controllers

import "github.com/revel/revel"

type ApiController struct {
  *revel.Controller
}

func (c ApiController) Index() revel.Result {
  return c.Render()
}
