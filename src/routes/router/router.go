package router

import (
	"net/http"
	"log"
	"runtime/debug"
)

//路由器接口，用于分发请求处理过程
type Router interface {

	/*
	 * 路由分发列表
	 */
	Routes() map[string]http.HandlerFunc
}

/*
 * 分发路由
 * @param serveMux *http.ServeMux
 * @param path 上层路由地址
 * @param router 路由器对象
 */
func Route(serveMux *http.ServeMux, path string, router Router) error {
	var routes map[string]http.HandlerFunc
	routes = router.Routes()

	if len(routes) > 0 {
		for name := range routes {
			uri := name
			if path != "/" {
				uri = path + name
			}
			serveMux.HandleFunc(uri, safeHandler(routes[name]))
		}
	}

	return nil
}

/*
 * 全局路由错误错误处理函数
 */
func safeHandler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request){
		defer func(){
			if e, ok := recover().(error); ok {
				http.Error(respWriter, e.Error(), http.StatusInternalServerError)

				//输出自定义的错误页面
				//respWriter.WriteHeader(http.StatusInternalServerError)
				//renderHtml(respWriter, "error", e)

				//日志
				log.Println("WARN: panic in %v. - %v", handlerFunc, e)
				log.Println(string(debug.Stack()))
			}
		}()
		handlerFunc(respWriter, request)
	}
}