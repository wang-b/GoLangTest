package routes

import (
	"net/http"
	"io/ioutil"
	"../config"
	"../renderer"
	"strings"
	"./router"
	. "../common"
)

type IndexRouter struct {

}

func (r *IndexRouter) Routes() map[string]router.RouterHandler {
	routes := make(map[string]router.RouterHandler)
	routes["/"] = r.indexHandler()
	return  routes
}

func (r *IndexRouter) Before(respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {

}

func (r *IndexRouter) After(respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {

}

func (r *IndexRouter) indexHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		fileInfos, err := ioutil.ReadDir(config.UPLOAD_DIR)
		CheckError(err)

		data := make(map[string]interface{})
		images := []string{}
		for _, fileInfo := range fileInfos {
			if strings.EqualFold(fileInfo.Name(), ".gitkeep") {
				continue
			}
			images = append(images, fileInfo.Name())
		}
		data["images"] = images
		renderer.RenderHtml(respWriter, "list", data)
	}
}
