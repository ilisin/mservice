package main

import (
	"github.com/Sirupsen/logrus"
	"gogs.xlh/bebp/mservice"
)

type API struct {
}

type InParam struct {
	Size int `json:"size"`
	Skip int `json:"skip"`
}

func (api *API) MSHello() (mservice.HTTPMethod, string, func(ctx *mservice.Context, p *InParam) (string, error)) {
	return mservice.GET, "/hello", func(ctx *mservice.Context, p *InParam) (string, error) {
		logrus.Info(p)
		logrus.Debug(ctx.QueryParam("p"))
		return "hello someone", nil
	}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	ms := mservice.NewMService()
	ms.AddPrototype("", &API{})
	ms.Run()
}
