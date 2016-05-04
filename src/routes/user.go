package routes

import (
	"net/http"
	"../renderer"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
	"encoding/json"
	"./router"
	. "../common"
)

const (
	DB_USER = "root"
	DB_PASSWORD = "123456"
	DB_NAME = "decade"
	DB_HOST = "127.0.0.1:3306"
)

const (
	INSERT_USER = "INSERT INTO user(username, password, nickname, realname, email, mobile, createTime) VALUES(?, ?, ?, ?, ?, ?, ?)"
	QUERY_USER = "SELECT * FROM user"
	UPDATE_USER = "UPDATE user SET username=?, password=?, nickname=?, realname=?, email=?, mobile=? WHERE id=?"
	DELETE_USER = "DELETE FROM user WHERE id=?"
	FIND_USER = "SELECT id, username, password, nickname, realname, email, mobile, createTime FROM user WHERE id=?"
)

type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Email string `json:"email"`
	Mobile string `json:"mobile"`
	CreateTime int64 `json:"createTime"`
}

func NewUser(id int64, username, password, nickname, realname, email, mobile string) *User {
	return &User{
		Id: id,
		Username: username,
		Password: password,
		Nickname: nickname,
		Realname: realname,
		Email: email,
		Mobile: mobile }
}

type UserRouter struct {

}

func (t *UserRouter) Routes() map[string]router.RouterHandler {
	routes := make(map[string]router.RouterHandler);
	routes["/"] = t.listHandler();
	routes["/add"] = t.addHandler();
	routes["/view"] = t.viewHandler();
	routes["/update"] = t.updateHandler();
	routes["/delete"] = t.deleteHandler();
	return routes;
}

func (r *UserRouter) Before(respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
	log.Println("userRouter.Before()...")
}

func (r *UserRouter) After(respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
	log.Println("userRouter.After()...")
}

func (t *UserRouter) listHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		db, err := openDB();
		defer db.Close()
		CheckError(err)

		rows, err := db.Query(QUERY_USER);
		defer rows.Close()
		CheckError(err)

		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i, _ := range values {
			scanArgs[i] = &values[i]
		}

		data := make(map[string]interface{})
		users := make([]*User, 0)
		var k int = 0

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			CheckError(err)
			record := make(map[string]string)
			for i, col := range values {
				if col != nil {
					record[columns[i]] = string(col.([]byte))  //col.([]byte)类型检查相关
				}
			}
			log.Println("record: ", record)

			//构建User对象
			user := &User{}
			user.Id, _ = strconv.ParseInt(record["id"], 10, 0)
			user.Username = record["username"]
			user.Password = record["password"]
			user.Nickname = record["nickname"]
			user.Realname = record["realname"]
			user.Email = record["email"]
			user.Mobile = record["mobile"]
			user.CreateTime, _ = strconv.ParseInt(record["createTime"], 10, 0)
			users = append(users, user)
			k ++
			log.Println("user: ", user)
		}

		data["users"] = users
		renderer.RenderHtml(respWriter, "user/list", data)
	}
}

func (t *UserRouter) addHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		if request.Method == "GET" {
			renderer.RenderHtml(respWriter, "user/add", nil)
 		} else if request.Method == "POST" {
			username := request.FormValue("username")
			password := request.FormValue("password")
			nickname := request.FormValue("nickname")
			realname := request.FormValue("realname")
			email := request.FormValue("email")
			mobile := request.FormValue("mobile")
			createTime := time.Now().Unix()

			db, err := openDB()
			defer db.Close()
			CheckError(err)

			stmt, err := db.Prepare(INSERT_USER)
			defer stmt.Close()
			CheckError(err)

			result, err := stmt.Exec(username, password, nickname, realname, email, mobile, createTime)
			CheckError(err)

			id, err := result.LastInsertId()
			CheckError(err)
			log.Println("insert user id: " + strconv.FormatInt(id, 10))

			http.Redirect(respWriter, request, "/user", http.StatusFound)
		}
	}
}

func (t *UserRouter) viewHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		id := request.FormValue("id")
		showType := request.FormValue("type")
		log.Println("请求参数: id = ", id, "showType = ", showType)
		if showType == "" {
			showType = "view"
		}

		userId, err := strconv.ParseInt(id, 10, 0)
		CheckError(err)

		user, err := findUserById(userId)
		CheckError(err)
		log.Println("找到user： ", user)

		data := make(map[string]interface{})
		data["user"] = user
		data["type"] = showType
		renderer.RenderHtml(respWriter, "user/view", data)
	}
}

func (t *UserRouter) updateHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		id := request.FormValue("id")
		if request.Method == "GET" {
			http.Redirect(respWriter, request, "/user/view?id=" + id + "&type=" + "update", http.StatusFound)
		} else if request.Method == "POST" {
			username := request.FormValue("username")
			password := request.FormValue("password")
			nickname := request.FormValue("nickname")
			realname := request.FormValue("realname")
			email := request.FormValue("email")
			mobile := request.FormValue("mobile")

			db, err := openDB()
			defer db.Close()
			CheckError(err)

			stmt, err := db.Prepare(UPDATE_USER)
			defer stmt.Close()
			CheckError(err)

			result, err := stmt.Exec(username, password, nickname, realname, email, mobile, id)
			CheckError(err)
			num, err := result.RowsAffected()
			log.Println("更新记录数: ", num)

			http.Redirect(respWriter, request, "/user/view?id=" + id + "&type=" + "view", http.StatusFound)
		}
	}
}

func (t *UserRouter) deleteHandler() router.RouterHandler {
	return func (respWriter http.ResponseWriter, request *http.Request, context *router.RouterContext) {
		userId := request.FormValue("id")

		db, err := openDB()
		defer db.Close()
		CheckError(err)

		stmt, err := db.Prepare(DELETE_USER)
		defer stmt.Close()
		CheckError(err)

		result, err := stmt.Exec(userId)
		CheckError(err)
		num, err := result.RowsAffected()
		CheckError(err)

		log.Println("更新记录数: ", num)
		//http.Redirect(respWriter, request, "/user", http.StatusFound)

		respWriter.Header().Add("Content-Type", "application/json")
		var jsonResult struct{   //匿名结构体
			Code int
			Result interface{}
			Message string
		}
		jsonResult.Code = 0
		jsonResult.Message = "删除成功"
		jsonData, err := json.Marshal(jsonResult)
		CheckError(err)
		respWriter.Write(jsonData)
	}
}

func openDB() (*sql.DB, error) {
	dbUrl := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dbUrl)
	return db, err
}

func findUserById(id int64) (*User, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(FIND_USER, id)
	log.Println("找到row: ", row)
	user := new(User)
	//注：此处查询需要值一一对应
	row.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.Realname, &user.Email, &user.Mobile, &user.CreateTime)
	return user, nil
}