package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koyeo/buck"
	"github.com/koyeo/buck/example/service"
	"github.com/koyeo/buck/exporter"
)

func main() {
	engine := gin.Default()
	engine.Use(cors.Default())
	
	api := buck.New()
	api.SetVersion("1.0.0")
	api.AddRouter(service.NewUserRouter(new(service.UserImplService)))
	api.SetExporter(":8099", &exporter.Options{
		Envs: []*exporter.Env{
			{
				Name: "本地测试",
				Host: "http://localhost:8088",
			},
		},
	})
	api.SetEngine(engine)
	api.Run(":8088")
}
