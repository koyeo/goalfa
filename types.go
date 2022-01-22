package _api

import (
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

type Register interface {
	//CONNECT(path string, handles ...gin.HandlerFunc) gin.IRoutes
	//TRACE(path string, handles ...gin.HandlerFunc) gin.IRoutes
	DELETE(path string, handles ...gin.HandlerFunc) gin.IRoutes
	GET(path string, handles ...gin.HandlerFunc) gin.IRoutes
	HEAD(path string, handles ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PATCH(path string, handles ...gin.HandlerFunc) gin.IRoutes
	POST(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PUT(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Any(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Group(prefix string, middleware ...gin.HandlerFunc) *gin.RouterGroup
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
}

type Route struct {
	Method      string
	Prefix      string
	Path        string
	Middlewares []gin.HandlerFunc `json:"-"`
	Children    []Route
	Handler     interface{} `json:"-"`
	handler     reflect.Value
}

type Router interface {
	Routes() []Route
}

type Module struct {
}

type Driver interface {
	Register(api *API) error
	Start(addr string) error
}

type HandlerInfo struct {
	Name     string
	Location string
}

type RouteTable struct {
	rows []RouteRow
}

func (p *RouteTable) AddRow(method, path string, info HandlerInfo) {
	p.rows = append(p.rows, RouteRow{
		Method:   method,
		Path:     path,
		Source:   info.Name,
		Location: info.Location,
	})
}

func (p RouteTable) Print() {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"", "接口地址", "方法名", "位置"})
	var data [][]string
	for _, v := range p.rows {
		data = append(data, []string{v.Method, v.Path, v.Source, v.Location})
	}
	t.AppendBulk(data)
	t.Render()
}

type RouteRow struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Source   string `json:"source"`
	Location string `json:"location"`
}
