package router

import (
	"net/http"
	"log"
	"runtime/debug"
)

/*
 * 路由对象上下文
 */
type RouterContext struct {

}

/*
 * 路由处理函数对象
 */
type RouterHandler func(http.ResponseWriter, *http.Request, *RouterContext)

//路由器接口，用于分发请求处理过程
type Router interface {

	/*
	 * 路由前置处理方法
	 */
	Before(http.ResponseWriter, *http.Request, *RouterContext)

	/*
	 * 路由后置处理函数
	 */
	After(http.ResponseWriter, *http.Request, *RouterContext)

	/*
	 * 路由分发列表
	 */
	Routes() map[string]RouterHandler
}

/*
 * 分发路由
 * @param serveMux *http.ServeMux
 * @param path 上层路由地址
 * @param router 路由器对象
 */
func Route(serveMux *http.ServeMux, path string, router Router) error {
	var routes map[string]RouterHandler
	routes = router.Routes()
	var context *RouterContext = new(RouterContext)  //创建上下文

	if len(routes) > 0 {
		for name := range routes {
			uri := name
			if path != "/" {
				uri = path + name
			}
			serveMux.HandleFunc(uri, routeHandler(router, context, routes[name]))
		}
	}

	return nil
}

/*
 * 调用路由处理函数，并处理异常
 * @param router 路由对象
 * @param context 上下文对象
 * @param handler 路由处理函数
 */
func routeHandler(router Router, context *RouterContext, handler RouterHandler) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request){
		defer func(){
			if e, ok := recover().(error); ok {
				http.Error(respWriter, e.Error(), http.StatusInternalServerError)

				//输出自定义的错误页面
				//respWriter.WriteHeader(http.StatusInternalServerError)
				//renderHtml(respWriter, "error", e)

				//日志
				log.Println("WARN: panic in %v. - %v", handler, e)
				log.Println(string(debug.Stack()))
			}
		}()
		router.Before(respWriter, request, context)
		handler(respWriter, request, context)
		router.After(respWriter, request, context)
	}
}