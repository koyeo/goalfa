package exporter

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gozelle/_log"
	"github.com/gozelle/_log/wrap"
	"github.com/koyeo/buck/assets"
	"github.com/koyeo/buck/utils"
	"github.com/ttacon/chalk"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func NewExporter(addr string, options *Options) *Exporter {
	e := &Exporter{addr: addr, options: options}
	e.initBasicTypes()
	return e
}

type Exporter struct {
	version    string
	addr       string
	options    *Options
	Name       string
	Package    string
	Methods    []*Method
	basicTypes map[string]bool
}

func (p *Exporter) Init(version string) {
	p.version = version
	p.initBasicTypes()
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
	data, err := sdk.Make(c.Query("lang"), c.Query("package"))
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

func (p *Exporter) initBasicTypes() {
	if p.options == nil {
		return
	}
	for _, v := range p.options.BasicTypes {
		if p.basicTypes == nil {
			p.basicTypes = map[string]bool{}
		}
		r := reflect.ValueOf(v.Elem)
		p.basicTypes[fmt.Sprintf("%s@%s", r.Type().PkgPath(), r.Type().String())] = true
	}
}

// ReflectFields 反射转换输入输出的字段信息
func (p Exporter) ReflectFields(name, label string, validator *Validator, t reflect.Type) (field *Field) {
	t = utils.TypeElem(t)
	field = new(Field)
	field.Name = name
	field.Label = label
	field.Type = p.getType(t)
	field.Validator = validator
	if t.Kind() == reflect.Struct && !p.isBasicType(t) {
		field.Struct = true
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			_name := p.getName(sf)
			_label := p.getFieldLabel(sf)
			_validator := p.getFieldValidator(sf)
			_field := p.ReflectFields(_name, _label, _validator, sf.Type)
			if _field.Struct || _field.Nested {
				field.Nested = true
			}
			field.Fields = append(field.Fields, _field)
		}
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		field.Array = true
		field.Elem = p.ReflectFields("", label, validator, t.Elem())
		if field.Elem.Struct || field.Elem.Nested {
			field.Nested = true
		}
	}
	return
}

func (p Exporter) isBasicType(t reflect.Type) bool {
	fmt.Println("nil", p.basicTypes)
	if p.basicTypes == nil {
		return false
	}
	fmt.Println(fmt.Sprintf("%s@%s", t.PkgPath(), t.String()))
	_, ok := p.basicTypes[fmt.Sprintf("%s@%s", t.PkgPath(), t.String())]
	return ok
}

func (p Exporter) getType(t reflect.Type) string {
	if p.isBasicType(t) {
		return t.String()
	}
	s := t.String()
	if strings.Contains(s, ".") {
		s = strings.Split(s, ".")[1]
	}
	return s
}

func (p Exporter) getFieldLabel(field reflect.StructField) string {
	return field.Tag.Get("label")
}

func (p Exporter) getFieldValidator(field reflect.StructField) (validator *Validator) {
	required := strings.Contains(field.Tag.Get("validator"), "required")
	if required {
		validator = p.newIfNoValidator(validator)
		validator.Required = true
	}
	return
}

func (p Exporter) newIfNoValidator(validator *Validator) *Validator {
	if validator == nil {
		validator = new(Validator)
	}
	return validator
}

func (p Exporter) getName(field reflect.StructField) string {
	n := field.Tag.Get("json")
	n = strings.ReplaceAll(n, ",omitempty", "")
	if n == "" {
		n = field.Name
	}
	return n
}
