# API 开发框架

## 特性

* 支持 Service 映射路由接口
* 方便实现 Mock 接口 
* 支持自动生成文档
* 自动自动生成 SDK

## 快速启动

```go
engine := gin.New() 

api := _api.New()
api.AddRouter(NewUserRouter())
api.SetVersion("1.0.0")
api.SetEngine(dirver)
api.SetExpoter(":8082")
api.Run(":8081")
```

## 服务定义

### 服务定义

```go
type UserService interface {
    Ping(ctx context.Context)(out string, err error)
    Register(ctx context.Context, in RegisterIn) error
    Login(ctx context.Context, in *LoginIn)(out LoginOut, err error)
}
```

### 方法格式

**格式说明:**
```
func (ctx context.Context [, in struct])([out interface{},] err error)
```

**入参说明:**

第一个参数必选，必须是 context.Context 类型，第二个参数可选，必须是一个 struct。

**出参说明:**

第一个参数 `out` 可选，可以为任意类型。第二个参数必选，必须是 error 类型，
如果 out 为 go 的基础数据类型，如 string、int、float64 等、或实现了 String() 方法，则 API 返回报文采用字符串编码，否则将采用 json 编码并输出。

## Router 定义

Router 用来描述服务的 API 匹配关系，同时可以在 Router 中定义中间件、前缀等。实现 `buck.Router` 即可实现一个 `Router`。

```go
type UserRouter struct {
	service UserService
}

func (u UserRouter) Routers() []buck.Route {
	return []buck.Route{
	    {
	    	Prefix: "/api",
	    	Middleware: SomeMiddleware,
	    	Children: []buck.Route{
	    	  { Method: buck.Get, Handler: p.service.Ping },
	    	  { Handler: p.service.Register },
	    	  { Path: "/login", Handler: p.service.Login },
            },
        }	
    }   
}
```