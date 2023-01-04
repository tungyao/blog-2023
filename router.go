package main

import (
	"crypto/sha1"
	"fmt"
	uc "github.com/tungyao/ultimate-cedar"
	"net/http"
	"time"
)

type Post struct {
	Name       string `json:"name"`
	Data       string `json:"data"`
	CreateTime int    `json:"create_time"`
	UpdateTIme int    `json:"update_t_ime"`
	Permission int    `json:"permission"`
	IsDelete   int    `json:"is_delete"`
}

type IndexCache struct {
}

func Index(writer uc.ResponseWriter, request uc.Request) {
	limit, offset := page(request)
	pageCache := Spruce.Get([]byte("index"))
	if pageCache == nil {
		Db.Query("select * from post limit ? offset ?", limit*offset, offset)
	}
	writer.Data(pageCache).Send()
}

func OnlyOne(writer uc.ResponseWriter, request uc.Request) {
	//data := Caches.Get(request.Data.Get("name"))
}

func MgPostGet(writer uc.ResponseWriter, request uc.Request) {

}

func MgPostAdd(writer uc.ResponseWriter, request uc.Request) {

}
func MgPostUpdate(writer uc.ResponseWriter, request uc.Request) {

}
func MgPostDelete(writer uc.ResponseWriter, request uc.Request) {

}
func FileUpload(writer uc.ResponseWriter, request uc.Request) {

}
func FileUpdate(writer uc.ResponseWriter, request uc.Request) {

}

type User struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

var user = make(map[string]*User)

// Login 信息保存在内存中
func Login(writer uc.ResponseWriter, request uc.Request) {
	u := User{}
	if err := IoRead(request, &u); err != nil {
		http.Error(writer, err.Error(), 403)
		return
	}
	row := Db.QueryRow("select count(*) from user where name=? and pass=?", u.Name, u.Pass)
	yes := 0
	row.Scan(&yes)
	if yes == 0 {
		http.Error(writer, "check failed", 401)
		return
	}
	sh := sha1.New()
	sh.Write([]byte(time.Now().String() + u.Pass + u.Name))
	token := fmt.Sprintf("%x", sh.Sum(nil))
	user[token] = &u
	writer.Data(fmt.Sprintf(`{"token":"%s"}`, token)).Send()
}
