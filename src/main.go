package main

import (
	"net/http"
	"log"
	"./system"
)

/*
 * 应用初始化函数，在main()之前执行
 */
func init() {

}

/*
 * 应用程序启动入口
 */
func main() {
	mux := system.NewServeMux()

	//开启服务器监听
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe： ", err.Error())
	} else {
		log.Println("ListenAndServe started!")
	}
}