package demo

import (
	"net/http"

	_ "github.com/NicholeGit/nade/app/provider/demo"
	"github.com/gin-gonic/gin"
)

type DemoApi struct {
	// service *Service
}

func NewDemoApi() *DemoApi {
	// service := NewService()
	// return &DemoApi{service: service}
	return &DemoApi{}
}

func Register(r *gin.Engine) error {
	api := NewDemoApi()
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
func (api *DemoApi) Demo(c *gin.Context) {
	c.JSON(http.StatusOK, "this is demo for dev all")
}
