package webservice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/knoguchi/go_project_template/services"
	"github.com/knoguchi/go_project_template/services/myservice"
	"net/http"
)

var registry *services.ServiceRegistry

func handleRoot(c *gin.Context) {
	var mysvc *myservice.MyService
	if err := registry.FetchService(&mysvc); err != nil {
		log.Errorf("%T is not found in registry.  Make sure it's enabled", mysvc)
		panic(err)
	}

	count := mysvc.GetCount()

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":  http.StatusOK,
			"error": "Welcome server 01 " + fmt.Sprintf("%d", count),
		},
	)
}
func router01(reg *services.ServiceRegistry) http.Handler {
	registry = reg
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", handleRoot)

	return e
}
