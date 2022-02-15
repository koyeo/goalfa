package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koyeo/buck"
	"github.com/koyeo/buck/example/foo/service"
	"github.com/koyeo/buck/exporter"
	"github.com/shopspring/decimal"
)

func main() {
	
	// gin 跨域配置
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowAllOrigins = true
	
	// 自定义 gin 驱动
	engine := gin.Default()
	engine.Use(cors.New(config))
	
	// buck 实例
	api := buck.New()
	api.SetVersion("1.0.0")
	api.SetEngine(engine)
	api.AddRouter(service.NewFooRouter(new(service.FooMock)))
	
	// API 导出器配置
	api.SetExporter(":9090", &exporter.Options{
		Project: "Foo",
		Envs: []exporter.Env{
			{
				Name: "本地测试",
				Host: "http://localhost:8080",
			},
		},
		BasicTypes: []exporter.BasicType{
			{
				Elem: decimal.Decimal{},
				Mapping: map[string]string{
					"ts": "string",
				},
			},
			{
				Elem: service.CID{},
				Mapping: map[string]string{
					"ts": "string",
				},
			},
		},
	})
	
	api.Run(":8080")
}
