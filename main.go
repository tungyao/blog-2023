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
	DbPath         = "./main.db"
	FilePathPrefix = "./file/"
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

	// 增加中间价
	login := uc.MiddlewareInterceptor(func(writer uc.ResponseWriter, request uc.Request, handlerFunc uc.HandlerFunc) {
		u := user[request.Header.Get("auth")]
		if u != nil {
			handlerFunc(writer, request)
		} else {
			writer.Data(`{"msg":"need login"}`).Status(401).Send()
			return
		}
	})
	middleware := uc.MiddlewareChain{
		login,
	}
	r := uc.NewRouter()

	// 路由器
	r.Get("/", Index)

	r.Get("page/:name", OnlyOne)

	r.Group("mg", func(groups *uc.Groups) {

		groups.Get("post", MgPostGet, middleware)

		groups.Post("post", MgPostAdd, middleware)

		groups.Put("post", MgPostUpdate, middleware)

		groups.Delete("post", MgPostDelete, middleware)

		groups.Post("file", FileUpload, middleware)

		groups.Put("file", FileUpdate, middleware)

		groups.Post("login", Login)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
