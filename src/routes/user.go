package routes

import (
	"net/http"
	"../renderer"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"./router"
	"log"
	"strconv"
	"time"
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
)

type User struct {
	Id int64
	Username string
	Password string
	Nickname string
	Realname string
	Email string
	Mobile string
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

func (t *UserRouter) Routes() map[string]http.HandlerFunc {
	routes := make(map[string]http.HandlerFunc);
	routes["/"] = t.listHandler();
	routes["/add"] = t.addHandler();
	routes["/update"] = t.updateHandler();
	routes["/delete"] = t.deleteHandler();
	return routes;
}

func (t *UserRouter) listHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		db := openDB();
		rows, err := db.Query(QUERY_USER);
		router.CheckError(err)

		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i, _ := range values {
			scanArgs[i] = &values[i]
		}

		data := make(map[string]interface{})
		users := make([]User, len(columns))

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			router.CheckError(err)
			record := make(map[string]string)
			for i, col := range values {
				if col != nil {
					record[columns[i]] = string(col.([]byte))  //col.([]byte)类型检查相关
				}
			}
			log.Println("record: ", record)
		}

		data["users"] = users
		renderer.RenderHtml(respWriter, "user/list", data)
	}
}

func (t *UserRouter) addHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			renderer.RenderHtml(respWriter, "user/view", nil)
 		} else if request.Method == "POST" {
			username := request.FormValue("username")
			password := request.FormValue("password")
			nickname := request.FormValue("nickname")
			realname := request.FormValue("realname")
			email := request.FormValue("email")
			mobile := request.FormValue("mobile")
			createTime := time.Now().Unix()

			db := openDB()
			stmt, err := db.Prepare(INSERT_USER)
			router.CheckError(err)
			result, err := stmt.Exec(username, password, nickname, realname, email, mobile, createTime)
			router.CheckError(err)
			id, err := result.LastInsertId()
			router.CheckError(err)
			log.Println("insert user id: " + strconv.FormatInt(id, 10))

			http.Redirect(respWriter, request, "/user", http.StatusFound)
		}
	}
}

func (t *UserRouter) updateHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {

	}
}

func (t *UserRouter) deleteHandler() http.HandlerFunc {
	return func (respWriter http.ResponseWriter, request *http.Request) {

	}
}

func openDB() (*sql.DB) {
	dbUrl := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dbUrl)
	router.CheckError(err)
	return db
}