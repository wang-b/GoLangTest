package routes

import (
	"net/http"
	"io/ioutil"
	"../config"
	"../renderer"
	"./router"
)

type IndexRouter struct {

}

func (r *IndexRouter) Routes() map[string]http.HandlerFunc {
	routes := make(map[string]http.HandlerFunc)
	routes["/"] = r.indexHandler()
	return  routes
}

func (r *IndexRouter) indexHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		fileInfos, err := ioutil.ReadDir(config.UPLOAD_DIR)
		router.CheckError(err)

		data := make(map[string]interface{})
		images := []string{}
		for _, fileInfo := range fileInfos {
			images = append(images, fileInfo.Name())
		}
		data["images"] = images
		renderer.RenderHtml(respWriter, "list", data)
	}
}