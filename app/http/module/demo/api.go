package demo

import (
	"net/http"
	// init
	_ "github.com/NicholeGit/nade/app/provider/demo"
	"github.com/gin-gonic/gin"
)

type DemoAPI struct {
	// service *Service
}

func NewDemoAPI() *DemoAPI {
	// service := NewService()
	// return &DemoAPI{service: service}
	return &DemoAPI{}
}

func Register(r *gin.Engine) error {
	api := NewDemoAPI()
	r.GET("/demo/demo", api.Demo)
	// r.GET("/demo/demo2", api.Demo2)
	// r.POST("/demo/demo_post", api.DemoPost)
	// r.GET("/demo/orm", api.DemoOrm)
	// r.GET("/demo/cache/redis", api.DemoRedis)
	return nil
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []string
// @Router /demo/demo [get]
func (api *DemoAPI) Demo(c *gin.Context) {
	c.JSON(http.StatusOK, "this is demo for dev all")
}
