package system

import (
	"net/http"
	"../config"
	"../libs/utils"
	"../routes/router"
	"../routes"
)

func NewServeMux() *http.ServeMux{
	mux := http.NewServeMux()
	staticResHandler(mux, config.STATIC_PREFIX, config.STATIC_DIR)
	dispatchRouter(mux)

	return mux
}

/*
 * 处理静态资源访问
 * @param mux *http.ServeMux
 * @param prefix 请求路径URI前缀
 * @param staticDir 静态资源目录
 */
func staticResHandler(mux *http.ServeMux, prefix string, staticDir string) {
	mux.HandleFunc(prefix, func(respWriter http.ResponseWriter, request *http.Request){
		file := staticDir + request.URL.Path[(len(prefix) - 1):]
		if exists := utils.IsFileExists(file); !exists {
			http.NotFound(respWriter, request)
			return
		}
		http.ServeFile(respWriter, request, file)
	})
}

/*
 * 分发路由器
 */
func dispatchRouter(serveMux *http.ServeMux) {
	var rIndex router.Router
	var rPicture router.Router
	var rUser router.Router
	rIndex = new(routes.IndexRouter)
	rPicture = new(routes.PictureRouter)
	rUser = new(routes.UserRouter)

	router.Route(serveMux, "/", rIndex)
	router.Route(serveMux, "/picture", rPicture)
	router.Route(serveMux, "/user", rUser)
}