package renderer

import (
	"html/template"
	"io/ioutil"
	"path"
	"log"
	"net/http"
	"strings"
	"errors"
	"../config"
)

//模板文件扩展名
const TEMPLATE_EXT = ".html"

//模板缓存
var templates = make(map[string]*template.Template)

/*
 * 包初始化函数，在导入包之前隐式执行
 */
func init() {
	//加载并缓存模板文件
	loadTemplates(templates, config.TEMPLATE_DIR, "", TEMPLATE_EXT)
}

/*
 * 递归装载对应目录下所有的模板文件
 * @param container 用于缓存模板a
 * @param templatePath 目录
 * @param prefix 模板名称前缀
 * @param ext 模板文件后缀名
 */
func loadTemplates(container map[string]*template.Template, templateDir string, prefix string, ext string) {
	log.Println("scan dir： " + templateDir)
	fileInfos, err := ioutil.ReadDir(templateDir)
	if err != nil {  //如果没有模板，直接返回
		return
	}

	var name, tmplPath, tmplName string
	for _, fileInfo := range fileInfos {
		name = fileInfo.Name()

		tmplPath = templateDir + "/" + name
		tmplName = name
		if prefix != "" {
			tmplName = prefix + "/" + name
		}

		if fileInfo.IsDir() {  //如果是文件夹，递归加载
			loadTemplates(container, tmplPath, tmplName, ext)
		} else {
			if suffix := path.Ext(name); suffix != ext {
				continue
			}
			log.Println("loading template： " + tmplName)
			tmpl := template.Must(template.ParseFiles(tmplPath))  //模板编译失败，会报错
			container[tmplName] = tmpl
		}
	}
}

/*
 * 渲染模板
 * @param respWriter 服务器输出对象
 * @param tmpl 模板名称
 * @param data 模板数据
 */
func RenderHtml(respWriter http.ResponseWriter, tmpl string, data map[string]interface{}) {
	if ok := strings.HasSuffix(tmpl, TEMPLATE_EXT); !ok {
		tmpl = tmpl + TEMPLATE_EXT
	}
	htmlTmpl, ok := templates[tmpl]
	if ok {
		err := htmlTmpl.Execute(respWriter, data)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("Template not found： " + tmpl))
	}
}