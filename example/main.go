package main

import (
	"github.com/gozelle/_api"
	"github.com/gozelle/_api/example/service"
)

func main() {
	api := _api.New()
	api.SetVersion("1.0.0")
	api.AddRouter(service.NewUserRouter(new(service.UserImplService)))
	api.SetExporter(":8099", nil)
	api.Run(":8088")
}
