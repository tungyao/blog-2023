package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	sp "github.com/tungyao/spruce-light"

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

	// Spruce 页面的缓存
	Spruce *sp.Hash
)

func main() {
	// 初始化数据库
	InitDb()

	// 又一个缓存器
	Caches = NewCache()

	// 缓存器
	Spruce = sp.CreateHash(1024)

	// 加载进缓存
	ReadFromDb()

	r := uc.NewRouter()

	// 路由器
	r.Get("/", Index)

	r.Get("page/:name", OnlyOne)

	r.Group("mg", func(groups *uc.Groups) {

		groups.Get("post", MgPostGet)

		groups.Post("post", MgPostAdd)

		groups.Put("post", MgPostUpdate)

		groups.Delete("post", MgPostDelete)

		groups.Post("file", FileUpload)

		groups.Put("file", FileUpdate)

		groups.Post("login", Login)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
