package _api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_api/exporter"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type String string
type HTML string

func New() *API {
	return &API{
		routeTable: &RouteTable{},
	}
}

type API struct {
	version    string
	routers    []Router
	engine     *gin.Engine
	routeTable *RouteTable
	exporter   *exporter.Exporter
	methods    []*exporter.Method
}

func (p *API) SetVersion(version string) {
	p.version = version
}

func (p *API) AddRouter(router ...Router) {
	p.routers = append(p.routers, router...)
}

func (p *API) SetEngine(engine *gin.Engine) {
	p.engine = engine
}

func (p *API) SetExporter(addr string, options *exporter.Options) {
	p.exporter = exporter.NewExporter(addr, options)
}

func (p *API) Run(addr string) {
	if p.engine == nil {
		p.engine = gin.Default()
	}
	var (
		routes []Route
		err    error
	)
	for _, router := range p.routers {
		routes, err = p.prepareRoutes(router.Routes())
		if err != nil {
			panic(err)
		}
		err = p.registerRoutes(p.engine, "", routes)
		if err != nil {
			panic(err)
			return
		}
	}
	p.routeTable.Print()
	if p.exporter != nil {
		p.makeExporter()
		p.exporter.Run()
	}
	err = p.engine.Run(addr)
	if err != nil {
		panic(err)
		return
	}
}

// 检查参数是否为 error 类型
func (p *API) isError(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

// 检查参数是否为 context 类型
func (p *API) isContext(v reflect.Type) bool {
	if v.Name() == "Context" && v.PkgPath() == "context" {
		return true
	}
	return false
}

// 检查参数是否接受的路由 Handler 格式
func (p *API) isHandler(t reflect.Type) error {
	if t.Kind() != reflect.Func {
		return fmt.Errorf("expect func")
	}
	if t.NumIn() != 1 && t.NumIn() != 2 {
		return fmt.Errorf("expect max 2 params")
	}
	if t.NumIn() == 2 {
		in := t.In(1)
		for {
			if in.Kind() != reflect.Ptr {
				break
			}
			in = in.Elem()
		}
		if in.Kind() != reflect.Struct {
			return fmt.Errorf("input only acept struct")
		}
	}
	if !p.isContext(t.In(0)) {
		return fmt.Errorf("expect context")
	}
	if t.NumOut() != 1 && t.NumOut() != 2 {
		return fmt.Errorf("expect max 2 output params")
	}
	if !p.isError(t.Out(t.NumOut() - 1)) {
		return fmt.Errorf("expect error")
	}
	return nil
}

// 反射路由 Handler, 并检查是否为可接受的格式
func (p *API) parseHandler(handler interface{}) (v reflect.Value, err error) {
	v = reflect.ValueOf(handler)
	if err = p.isHandler(v.Type()); err != nil {
		err = fmt.Errorf("unexpect handler: %s", v.Type())
		return
	}
	return
}

// 预处理路由，反射路由处理器，并检查类型
func (p *API) prepareRoutes(in []Route) (out []Route, err error) {
	out = make([]Route, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[i]
		if out[i].Handler != nil {
			out[i].handler, err = p.parseHandler(out[i].Handler)
			if err != nil {
				return
			}
			if out[i].Method == "" {
				out[i].Method = http.MethodPost
			}
		}
		out[i].Children, err = p.prepareRoutes(out[i].Children)
		if err != nil {
			return
		}
	}
	return
}

// 递归注册路由树，处理中间件前缀逻辑，代理路由处理器为 Gin 控制器
func (p *API) registerRoutes(register Register, prefix string, routes []Route) (err error) {
	for _, v := range routes {
		if !v.handler.IsValid() {
			err = p.registerRoutes(
				register.Group(v.Prefix, v.Middlewares...),
				strings.Join([]string{prefix, v.Prefix}, ""),
				v.Children,
			)
			if err != nil {
				return
			}
		} else {
			p.routeTable.AddRow(v.Method, strings.Join([]string{prefix, v.Path}, ""), p.parseHandlerInfo(v.Handler))
			p.addMethod(v.Method, strings.Join([]string{prefix, v.Path}, ""), v.handler)
			switch v.Method {
			case http.MethodGet:
				register.GET(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodPost:
				register.POST(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodPut:
				register.PUT(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodDelete:
				register.DELETE(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodHead:
				register.HEAD(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			case http.MethodOptions:
				register.OPTIONS(v.Path, append([]gin.HandlerFunc{p.proxyHandler(v.handler)}, v.Middlewares...)...)
			default:
				err = fmt.Errorf("unsupport method: %s", v.Method)
				return
			}
		}
	}
	return
}

func (p *API) proxyHandler(handler reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		var out []reflect.Value
		ctx := context.Background()
		if handler.Type().NumIn() == 2 {
			var in reflect.Value
			var err error
			in, err = bindJson(c, handler.Type().In(1))
			if err != nil {
				// TODO 响应 Error 报错
				c.Error(err)
				return
			}
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx), in})
		} else {
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
		
		l := len(out)
		if !out[l-1].IsNil() {
			c.JSON(http.StatusInternalServerError, &Status{
				Detail: out[l-1].Interface().(error).Error(),
			})
			return
		}
		if l == 2 {
			r := out[0].Interface()
			switch r.(type) {
			case String:
				c.String(http.StatusOK, string(r.(String)))
				return
			default:
				c.JSON(http.StatusOK, r)
				return
			}
		}
		return
	}
}

func realType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t
		}
		t = t.Elem()
	}
}

func bindJson(c *gin.Context, t reflect.Type) (reflect.Value, error) {
	ptr := t.Kind() == reflect.Ptr
	if ptr {
		t = realType(t)
	}
	in := reflect.New(t)
	err := c.Bind(in.Interface())
	if err != nil {
		return in, err
	}
	if ptr {
		return in, nil
	}
	return in.Elem(), nil
}

// 解析 Handler 的信息
func (p *API) parseHandlerInfo(h interface{}) HandlerInfo {
	target := reflect.ValueOf(h).Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

// 解析 Handler 的信息
func (p *API) parseHandlerInfoValue(v reflect.Value) HandlerInfo {
	target := v.Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

// 生成 Exporter 信息
func (p API) makeExporter() {
	if p.exporter == nil {
		return
	}
	p.exporter.Methods = p.methods
}

func (p *API) addMethod(method, path string, handler reflect.Value) {
	info := p.parseHandlerInfoValue(handler)
	m := &exporter.Method{
		Name:   info.Name,
		Path:   path,
		Method: method,
	}
	if handler.Type().NumIn() > 1 {
		m.Input = exporter.ReflectFields(fmt.Sprintf("%sIn", m.Name), handler.Type().In(1))
	}
	if handler.Type().NumOut() > 1 {
		m.Output = exporter.ReflectFields(fmt.Sprintf("%sOut", m.Name), handler.Type().Out(0))
	}
	p.methods = append(p.methods, m)
	
}
