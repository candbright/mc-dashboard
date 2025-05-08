package route

import (
	"net/http"

	"github.com/candbright/go-server/pkg/rest"
	"github.com/gin-gonic/gin"
)

func init() {
	registerRoute(func(e *gin.Engine) {
		e.POST("/server/:id/allowlist/get", rest.H(getAllowList))
		e.POST("/server/:id/allowlist/add", rest.H(addAllowList))
		e.POST("/server/:id/allowlist/delete", rest.H(deleteAllowList))
	})
}

type AllowListAddReq struct {
	Username string `json:"username"`
}

type AllowListDeleteReq struct {
	Username string `json:"username"`
}

func getAllowList(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	allowList, err := server.GetAllowList()
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	return rest.Json(allowList)
}

func addAllowList(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	var req AllowListAddReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusBadRequest, err)
	}
	err = server.AllowListAdd(req.Username)
	if err != nil {
		return err
	}
	return nil
}

func deleteAllowList(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	var req AllowListDeleteReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusBadRequest, err)
	}
	err = server.AllowListDelete(req.Username)
	if err != nil {
		return err
	}
	return nil
}
