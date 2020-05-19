package main

import (
	//"bytes"
	//"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"html/template"
	"log"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	//"net/http"
	"os"
	"path/filepath"
)

var basicAuthPrefix = []byte("Basic ")

// BasicAuth is the basic auth handler
/*func BasicAuth(h fasthttprouter.Handle, user, pass[] byte) fasthttprouter.Handle {
	return fasthttprouter.Handle( func(ctx * fasthttp.RequestCtx, ps fasthttprouter.Params) {
		// Get the Basic Authentication credentials
	auth : =  ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix( auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], pass) {
					// Delegate request to the given handle
					h(ctx, ps)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	})
}*/

var Db *sqlx.DB

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Not protected!\n")
}

// Protected is the Protected handler
/*func Protected(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	fmt.Fprint(ctx, "Protected!\n")
}*/

func main() {

	println(os.Getenv("DB_DIALECT"))
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dialect := os.Getenv("DB_DIALECT")
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, dbHost, dbName)

	db, err := gorm.Open(dialect, connString)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Vse ogon")

	//создаем и запускаем в работу роутер для обслуживания запросов
	/*r := fasthttprouter.New()
	routes(r)
	//прикрепляемся к хосту и свободному порту для приема и обслуживания входящих запросов
	//вторым параметром передается роутер, который будет работать с запросами
	er := fasthttp.ListenAndServe("localhost:4444", r)
	if er != nil {
		log.Fatal(er)
	}*/
	//user := []byte("gordon")
	//pass := []byte("secret!")

	router := fasthttprouter.New()
	router.GET("/", Index)
	//router.GET("/protected/", BasicAuth(Protected, user, pass))

	log.Fatal(fasthttp.ListenAndServe("localhost:4444", router.Handler))

}

func StartPage(rw fasthttp.ResponseWriter, r *fasthttp.RequestCtx, p fasthttprouter.Params) {

	//указываем путь к нужному файлу
	path := filepath.Join("public", "html", "startStaticPage.html")

	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func routes(r *fasthttprouter.Router) {
	//путь к папке со внешними файлами: html, js, css, изображения и т.д.
	r.ServeFiles("/public/*filepath", fasthttp.Dir("public"))
	//что следует выполнять при входящих запросах указанного типа и по указанному адресу
	r.GET("/", StartPage)

}
