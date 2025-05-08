package route

import "github.com/gin-gonic/gin"

var routeIncubators []func(*gin.Engine)

func registerRoute(f func(*gin.Engine)) {
	routeIncubators = append(routeIncubators, f)
}

func Incubate(engine *gin.Engine) {
	for _, routeIncubator := range routeIncubators {
		routeIncubator(engine)
	}
}
