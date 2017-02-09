package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/mservice"
)

type API struct {
}

type InParam struct {
	Size int `req:"'size' min(20)"`
	Skip int `req:"enum(1;2;3)"`
}

func (api *API) APIDescription() map[string]string {
	return map[string]string{
		"": "",
	}
}

func (api *API) MSHello() (mservice.HTTPMethod, string, string, func(ctx *mservice.Context, p *InParam) (string, error)) {
	return mservice.GET,
		"/hello",
		"测试接口",
		func(ctx *mservice.Context, p *InParam) (string, error) {
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
