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

//模板缓存
var templates = make(map[string]*template.Template)

/*
 * 包初始化函数，在导入包之前隐式执行
 */
func init() {
	//加载并缓存模板文件
	fileInfos, err := ioutil.ReadDir(config.TEMPLATE_DIR)
	if err != nil {
		log.Println("templates not found...")
		return
	}

	var templateName, templatePath string
	for _, fileInfo := range fileInfos {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = config.TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template: ", templatePath)
		tmpl := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = tmpl
	}
}

/*
 * 渲染模板
 * @param respWriter 服务器输出对象
 * @param tmpl 模板名称
 * @param data 模板数据
 */
func RenderHtml(respWriter http.ResponseWriter, tmpl string, data map[string]interface{}) {
	if ok := strings.HasSuffix(tmpl, ".html"); !ok {
		tmpl = tmpl + ".html"
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