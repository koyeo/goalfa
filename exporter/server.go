package exporter

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_log"
	"github.com/gozelle/_log/wrap"
	"time"
)

func NewExporter(addr string, options *Options) *Exporter {
	return &Exporter{addr: addr, options: options}
}

type Exporter struct {
	addr    string
	options *Options
	Name    string
	Methods []*Method
}

func (p Exporter) Run() {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))
	// 导出SDK
	engine.GET("/sdk", func(c *gin.Context) {
		c.String(200, "some SDK")
	})
	// 导出描述协议
	engine.GET("/protocol", func(c *gin.Context) {
		c.JSON(200, p.Methods)
	})
	go func() {
		err := engine.Run(p.addr)
		if err != nil {
			_log.Panic("接口导出器启动失败", wrap.Error(err))
		}
	}()
}
