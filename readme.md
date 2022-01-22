# API 开发框架

## 特性

* 支持 Service 映射路由接口
* 方便实现 Mock 接口 
* 支持自动生成文档
* 自动自动生成 SDK

## 草稿

```go
engine := gin.New() 

api := _api.New()
api.AddRouter(NewUserRouter())
api.SetVersion("1.0.0")
api.SetEngine(dirver)
api.SetExpoter(":8082")
api.Run(":8081")
```