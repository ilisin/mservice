mservice
========

#### Rapid development framework of micro service API based on echo

## Feature Overview

- Modularization
- Middleware
- Swagger support 
- Session
- Authentication and authorization

## Quick Start

### Installation

**Windows**
    
    git clone http://gogs.xlh/bebp/mservice %GOPATH%\\src\\bebp\\mservice

**Mac OX or Linux**

    git clone http://gogs.xlh/bebp/mservice $GOPATH/src/bebp/mservice
    
    
### Hello, World!


    type API struct {
    }
    
    func (api *API) MSHello() (mservice.HTTPMethod, string, mservice.MSHandlerFunc) {
    	return mservice.GET, "/hello", func(ctx *mservice.Context) error {
    		return ctx.Context.String(http.StatusOK, "hello world!")
    	}
    }
    
    func main() {
    	logrus.SetLevel(logrus.DebugLevel)
    	ms := mservice.NewMService()
    	ms.AddPrototype("", &API{})
    	ms.Run()
    }

**Run**

    # curl -X GET http://localhost:8000/hello
    hello world!
