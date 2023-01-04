package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	uc "github.com/tungyao/ultimate-cedar"
	"log"
	"net/http"
)

const (
	// DbPath sqlite 数据库地址
	DbPath = "./main.db"
)

var (
	Db     *sql.DB
	Caches *Cache
)

func main() {
	// 初始化各种东西
	InitDb()
	Caches = NewCache()

	r := uc.NewRouter()

	// 路由器
	r.Get("page/:name", Index)

	r.Group("mg", func(groups *uc.Groups) {

		groups.Get("post", func(writer uc.ResponseWriter, request uc.Request) {

		})
		groups.Post("post", func(writer uc.ResponseWriter, request uc.Request) {

		})
		groups.Put("post", func(writer uc.ResponseWriter, request uc.Request) {

		})
		groups.Delete("post", func(writer uc.ResponseWriter, request uc.Request) {

		})
		groups.Post("file", func(writer uc.ResponseWriter, request uc.Request) {

		})
		groups.Delete("file", func(writer uc.ResponseWriter, request uc.Request) {

		})
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
