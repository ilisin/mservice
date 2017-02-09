package mservice

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type MService struct {
	e      *echo.Echo
	config *Config
	auth   Auther
	models []*Model
}

func (ms *MService) Route() {
	for _, m := range ms.models {
		for _, meta := range m.ReadHandlers() {
			var echoHandler func(string, echo.HandlerFunc, ...echo.MiddlewareFunc)
			switch meta.Method {
			case CONNECT:
				echoHandler = ms.e.CONNECT
			case DELETE:
				echoHandler = ms.e.DELETE
			case GET:
				echoHandler = ms.e.GET
			case HEAD:
				echoHandler = ms.e.HEAD
			case OPTIONS:
				echoHandler = ms.e.OPTIONS
			case PATCH:
				echoHandler = ms.e.PATCH
			case POST:
				echoHandler = ms.e.POST
			case PUT:
				echoHandler = ms.e.PUT
			case TRACE:
				echoHandler = ms.e.TRACE
			default:
				echoHandler = ms.e.GET
			}
			echoHandler(meta.Path, func(c echo.Context) error {
				context := &Context{
					Context: c,
				}
				return meta.Handler(context)
			})
		}
	}
}

func (ms *MService) Middleware() {
	ms.e.Use(middleware.Logger())
}

// add api prototype
func (ms *MService) AddPrototype(prefix string, p interface{}) {
	ms.models = append(ms.models, NewModel(prefix, p))
}

// server start
func (ms *MService) Run() {
	ms.Middleware()
	ms.Route()
	ms.e.Logger.Infof("Service run at %v", ms.config.Host)
	ms.e.Start(ms.config.Host)
}

func NewMService() *MService {
	ms := &MService{
		models: make([]*Model, 0),
	}
	ms.e = echo.New()
	ms.e.Logger.SetLevel(log.DEBUG)
	config, err := LoadAConfig()
	if err != nil {
		ms.e.Logger.Fatal(err)
	}
	ms.config = config
	return ms
}
