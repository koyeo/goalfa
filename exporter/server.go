package exporter

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_log"
	"github.com/gozelle/_log/wrap"
	"github.com/koyeo/buck/assets"
	"github.com/ttacon/chalk"
	"net/http"
	"strings"
	"time"
)

func NewExporter(addr string, options *Options) *Exporter {
	return &Exporter{addr: addr, options: options}
}

type Exporter struct {
	version string
	addr    string
	options *Options
	Name    string
	Methods []*Method
}

func (p *Exporter) SetVersion(version string) {
	p.version = version
}

func (p Exporter) Run() {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.GET("/sdk", p.sdkHandler)
	engine.GET("/protocol", p.protocolHandler)
	// engine.GET("/struct", p.protocolHandler)
	engine.StaticFS("/exporter", assets.Root)
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/exporter/index.html")
	})
	engine.GET("/:path", func(c *gin.Context) {
		path := c.Param("path")
		if path == "exporter" {
			path = "index.html"
		}
		if !strings.HasPrefix(path, "exporter") {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/exporter/%s", path))
		}
	})
	engine.GET("/test", func(c *gin.Context) {
		c.Request.URL.Path = "/sdk"
		engine.HandleContext(c)
	})
	go func() {
		p.printAddress()
		err := engine.Run(p.addr)
		if err != nil {
			_log.Panic("接口导出器启动失败", wrap.Error(err))
		}
	}()
}

// 打印 API 调试器访问地址
func (p Exporter) printAddress() {
	addr := p.addr
	if strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf("http://127.0.0.1%s", addr)
	} else {
		addr = fmt.Sprintf("http://%s", addr)
	}
	fmt.Println(chalk.Green.Color(strings.Repeat("=", 100)))
	fmt.Println(chalk.Green.Color(fmt.Sprintf("API 调试器访问地址：%s", addr)))
	fmt.Println(chalk.Green.Color(strings.Repeat("=", 100)))
}

// 导出 SDK 代码
func (p Exporter) sdkHandler(c *gin.Context) {
	sdk := NewSDK(p.Methods)
	data, err := sdk.Make(c.Query("lang"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

type ProtocolOutput struct {
	Version string    `json:"version"`
	Options *Options  `json:"options"`
	Methods []*Method `json:"methods"`
}

// 导出接口描述协议
func (p Exporter) protocolHandler(c *gin.Context) {
	
	out := new(ProtocolOutput)
	out.Version = p.version
	out.Options = p.options
	out.Methods = p.convertMethodTypes(c.Query("lang"))
	c.JSON(200, out)
}

func (p Exporter) convertMethodTypes(lang string) []*Method {
	methods := make([]*Method, 0)
	switch lang {
	case "ts":
		for _, v := range p.Methods {
			n := v.Fork()
			p.toTypescriptFieldType(n.Input)
			p.toTypescriptFieldType(n.Output)
			methods = append(methods, n)
		}
	default:
		methods = p.Methods
	}
	return methods
}

func (p Exporter) toTypescriptFieldType(field *Field) {
	if field == nil {
		return
	}
	field.Origin = field.Type
	field.Type = typescriptTypeConverter(field.Type)
	for _, v := range field.Fields {
		p.toTypescriptFieldType(v)
	}
	if field.Elem != nil {
		p.toTypescriptFieldType(field.Elem)
	}
}
