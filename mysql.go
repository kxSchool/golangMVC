//数据库连接池测试
package main

import (
	"database/sql"
	"fmt"
	"log"

	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:ml123456@tcp(127.0.0.1:3306)/iissy?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

func main() {
	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/pool", pool)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pool(w http.ResponseWriter, r *http.Request) {
	var Subject string
	var Description string
	a := rand.Intn(9)
	rows, err := db.Query("SELECT Subject,Description FROM article where Id=?", a)
	defer rows.Close()
	checkErr(err)

	for rows.Next() {
		err := rows.Scan(&Subject, &Description)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s", Subject, Description)
	//fmt.Fprintf(w, "%s", string(b))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
