package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"golang.com/ado"
	"golang.com/models"
	"golang.com/utils"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", utils.SQLDB)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

// NewsIndex 图文首页
func NewsIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	temp, _ := template.ParseFiles(
		"public/views/index.html",
		"public/views/shared/_list.picture.html",
		"public/views/shared/_header.less.html",
		"public/views/shared/_toper.html",
		"public/views/shared/_footer.html")

	result, err := ado.NewsIndex(db)
	if err != nil {
		log.Fatal(err)
	}

	err = temp.Execute(w, result)
	if err != nil {
		fmt.Fprintf(w, "%q", err)
	}
}

// TechIndex 技术文档首页
func TechIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	temp, _ := template.ParseFiles(
		"public/views/tech.html",
		"public/views/shared/_list.article.html",
		"public/views/shared/_header.less.html",
		"public/views/shared/_toper.html",
		"public/views/shared/_footer.html")

	result, err := ado.TechIndex(db)
	if err != nil {
		log.Fatal(err)
	}

	err = temp.Execute(w, result)
	if err != nil {
		fmt.Fprintf(w, "%q", err)
	}
}

// Upload 上传
func Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := time.Now()
	dir := t.Format("20060102")
	var result models.Uploador
	err := r.ParseMultipartForm(32 << 20)
	utils.CheckErr(err)

	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	defer file.Close()
	ext := utils.ExtensionName(handler.Filename)
	fmt.Println(ext)
	path := "/upload/" + dir + "/"
	err = os.MkdirAll("./public"+path, os.ModePerm)
	utils.CheckErr(err)
	path += utils.UniqueID() + strings.ToLower(ext)
	f, err := os.OpenFile("./public"+path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	result.Success = true
	result.Path = path
	defer f.Close()
	io.Copy(f, file)

	b, _ := json.Marshal(result)
	fmt.Fprintf(w, "%s", string(b))
}
