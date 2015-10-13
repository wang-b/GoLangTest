package routes

import (
	"net/http"
	"log"
	"strconv"
	"time"
	"os"
	"io"
	"../renderer"
	"../libs/utils"
	"../config"
	"./router"
	"strings"
)

type MainRouter struct {

}

func (r *MainRouter) Routes() map[string]http.HandlerFunc {
	routes := make(map[string]http.HandlerFunc)
	routes["/upload"] = r.uploadHandler()
	routes["/view"] = r.viewHandler()
	routes["/delete"] = r.deleteHandler()
	return  routes
}

func (r *MainRouter) uploadHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			renderer.RenderHtml(respWriter, "upload", nil)
		}
		if request.Method == "POST" {
			file, header, err := request.FormFile("image")
			router.CheckError(err)

			filename := header.Filename
			log.Println("upload file： " + strings.TrimSpace(filename))
			defer file.Close()

			fName := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filename

			//temp, err := ioutil.TempFile(UPLOAD_DIR, fName)  //此方法创建文件文件名有后缀
			temp, err := os.Create(config.UPLOAD_DIR + "/" + fName)
			router.CheckError(err)
			defer temp.Close()

			_, e := io.Copy(temp, file)
			router.CheckError(e)

			http.Redirect(respWriter, request, "/" + fName, http.StatusFound)
		}
	}
}

func (r *MainRouter) viewHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		imageId := request.FormValue("id")
		imagePath := config.UPLOAD_DIR + "/" + imageId
		if exists := utils.IsFileExists(imagePath); !exists {
			http.NotFound(respWriter, request)
			return
		}

		respWriter.Header().Set("Content-Type", "image")
		http.ServeFile(respWriter, request, imagePath)
	}
}

func (r *MainRouter)deleteHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		imageId := request.FormValue("id");
		imagePath := config.UPLOAD_DIR + "/" + imageId
		if exists := utils.IsFileExists(imagePath); !exists {
			http.NotFound(respWriter, request)
			return
		}
		log.Println("imagePath: " + imagePath)
		err := os.Remove(imagePath)
		if err != nil {
			http.Error(respWriter, err.Error(), http.StatusForbidden)
			return
		}

		//刷新页面
		http.Redirect(respWriter, request,"/", http.StatusFound)
	}
}
